from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework.viewsets import ModelViewSet
from django.http import Http404
from django.core.cache import cache
from django.db.models import Q
from data.models import Book, KeywordsEnglish, KeywordsFrench, KeywordBookEnglish, KeywordBookFrench
from data.serializers import BookSerializer
from data.jaccard import jaccard_similarity
import numpy as np
from collections import defaultdict
import time

class BookViewSet(APIView):
    
    def get(self, request, format=None):
        start_time = time.time()
        
        # Call the processing function
        queryset = self.process_book_query(request)
        
        # Serialize and return the response
        serializer = BookSerializer(queryset, many=True)
        
        end_time = time.time()
        execution_time = end_time - start_time
        print(f"BookViewSet query execution time: {execution_time:.4f} seconds")
        
        return Response(serializer.data)
    
    def process_book_query(self, request):
        # Initialize base queryset
        queryset = Book.objects.exclude(download_count__isnull=True)
        queryset = queryset.exclude(title__isnull=True)
        
        # Process language filter
        language = request.GET.get('languages')
        if language is not None:
            queryset = queryset.filter(languages__code=language)
        
        # Process author name filter
        queryset = self._filter_by_author(request, queryset)
        
        # Process title filter
        queryset = self._filter_by_title(request, queryset)
        
        # Process keyword filter
        queryset = self._filter_by_keyword(request, queryset, language)
        
        # Apply sorting
        queryset = self._apply_sorting(request, queryset)
        
        return queryset.distinct()
    
    def _filter_by_author(self, request, queryset):
        search_name_author = request.GET.get('author_name')
        if search_name_author is not None:
            search_name_authors_type = request.GET.get('author_name_type')
            search_name_authors_type = "classique" if search_name_authors_type is None else search_name_authors_type
            
            if search_name_authors_type == "classique":
                queryset = queryset.filter(authors__name__icontains=search_name_author)
            else:
                queryset = queryset.filter(authors__name__regex=search_name_author)
        return queryset
    
    def _filter_by_title(self, request, queryset):
        search_title = request.GET.get('title')
        if search_title is not None:
            search_title_type = request.GET.get('title_type')
            search_title_type = "classique" if search_title_type is None else search_title_type
            
            if search_title_type == "classique":
                queryset = queryset.filter(title__icontains=search_title)
            else:
                queryset = queryset.filter(title__regex=search_title)
        return queryset
    
    def _filter_by_keyword(self, request, queryset, language):
        search_keyword = request.GET.get('keyword')
        if search_keyword is not None:
            search_keywords_type = request.GET.get('keyword_type')
            search_method = 'icontains' if search_keywords_type == 'classique' else 'regex'
            
            # Define language-specific filters
            filters = {
                'en': {'keywordbookenglish__keyword__token__{}'.format(search_method): search_keyword},
                'fr': {'keywordbookfrench__keyword__token__{}'.format(search_method): search_keyword},
            }
            
            # Apply language-specific filter or both if language is not specified
            if language in filters:
                queryset = queryset.filter(**filters[language])
            else:
                # For other languages or no language specified, search in both English and French
                english_filter = Q(**{'keywordbookenglish__keyword__token__{}'.format(search_method): search_keyword})
                french_filter = Q(**{'keywordbookfrench__keyword__token__{}'.format(search_method): search_keyword})
                queryset = queryset.filter(english_filter | french_filter)
        return queryset
    
    def _apply_sorting(self, request, queryset):
        sort = request.GET.get('sort')
        if sort == 'download_count':
            ord = request.GET.get('order')
            ord = "descending" if ord is None else ord
            if ord == "descending":
                queryset = queryset.order_by('-download_count')
            else:
                queryset = queryset.order_by('download_count')
        return queryset


class BooksList(APIView):
    def get(self, request, format=None):
        start_time = time.time()
        
        # Use BookViewSet's functionality to process the query
        book_view = BookViewSet()
        queryset = book_view.process_book_query(request)
        
        # Serialize the queryset
        serializer = BookSerializer(queryset, many=True)
        results = serializer.data
        
        # Generate suggestions from the first result using Jaccard similarity
        suggestions = []
        if results and len(results) > 0:
            first_book_id = results[0]['id']  # gutenberg_id
            
            # Get language filter for consistency
            language = request.GET.get('languages', 'en')
            search_keyword = request.GET.get('keyword', '')
            
            print(f"[JACCARD SUGGESTIONS] First book ID: {first_book_id}, Keyword: {search_keyword}")
            
            # If we have a keyword and first result, get similar books using Jaccard similarity
            if search_keyword:
                try:
                    # Build keyword occurrence dict for first book
                    lang_mapping = {'en': 'english', 'fr': 'french'}
                    search_language = lang_mapping.get(language, 'both')
                    
                    # Get keywords and their occurrences for first book
                    first_book_keywords = {}
                    
                    if search_language in ['english', 'both']:
                        english_kws = KeywordBookEnglish.objects.filter(
                            book_id=first_book_id
                        ).select_related('keyword')
                        for kw_book in english_kws:
                            first_book_keywords[kw_book.keyword.token] = kw_book.occurence
                    
                    if search_language in ['french', 'both']:
                        french_kws = KeywordBookFrench.objects.filter(
                            book_id=first_book_id
                        ).select_related('keyword')
                        for kw_book in french_kws:
                            first_book_keywords[kw_book.keyword.token] = kw_book.occurence
                    
                    print(f"[JACCARD SUGGESTIONS] First book has {len(first_book_keywords)} keywords")
                    
                    if first_book_keywords:
                        # Get all books with any of these keywords
                        similar_books_with_scores = []
                        
                        if search_language in ['english', 'both']:
                            keyword_objs = KeywordsEnglish.objects.filter(
                                token__in=first_book_keywords.keys()
                            )
                            book_ids = KeywordBookEnglish.objects.filter(
                                keyword__in=keyword_objs
                            ).exclude(book_id=first_book_id).values_list('book_id', flat=True).distinct()
                        else:
                            book_ids = []
                        
                        if search_language in ['french', 'both']:
                            keyword_objs_fr = KeywordsFrench.objects.filter(
                                token__in=first_book_keywords.keys()
                            )
                            book_ids_fr = KeywordBookFrench.objects.filter(
                                keyword__in=keyword_objs_fr
                            ).exclude(book_id=first_book_id).values_list('book_id', flat=True).distinct()
                            book_ids = set(book_ids) | set(book_ids_fr)
                        
                        print(f"[JACCARD SUGGESTIONS] Found {len(book_ids)} candidate books")
                        
                        # Calculate Jaccard similarity for each candidate book
                        for book_id in list(book_ids)[:100]:  # Limit to first 100 candidates
                            candidate_keywords = {}
                            
                            if search_language in ['english', 'both']:
                                cand_kws = KeywordBookEnglish.objects.filter(
                                    book_id=book_id
                                ).select_related('keyword')
                                for kw_book in cand_kws:
                                    candidate_keywords[kw_book.keyword.token] = kw_book.occurence
                            
                            if search_language in ['french', 'both']:
                                cand_kws_fr = KeywordBookFrench.objects.filter(
                                    book_id=book_id
                                ).select_related('keyword')
                                for kw_book in cand_kws_fr:
                                    candidate_keywords[kw_book.keyword.token] = kw_book.occurence
                            
                            # Calculate Jaccard similarity
                            if candidate_keywords:
                                similarity = jaccard_similarity(first_book_keywords, candidate_keywords)
                                similar_books_with_scores.append((book_id, similarity))
                        
                        # Sort by similarity score (descending) and take top 10
                        similar_books_with_scores.sort(key=lambda x: x[1], reverse=True)
                        top_similar_ids = [book_id for book_id, score in similar_books_with_scores[:10]]
                        
                        print(f"[JACCARD SUGGESTIONS] Top similar books: {len(top_similar_ids)}")
                        
                        if top_similar_ids:
                            # Fetch and serialize the books
                            suggestion_books = Book.objects.filter(gutenberg_id__in=top_similar_ids)
                            suggestions = BookSerializer(suggestion_books, many=True).data
                            print(f"[JACCARD SUGGESTIONS] Returning {len(suggestions)} suggestions")
                    else:
                        print(f"[JACCARD SUGGESTIONS] First book has no keywords")
                except Exception as e:
                    print(f"[JACCARD SUGGESTIONS] Error: {e}")
                    import traceback
                    traceback.print_exc()
                    suggestions = []
            else:
                print(f"[JACCARD SUGGESTIONS] No keyword provided")
        
        execution_time = time.time() - start_time
        print(f"BookList query execution time: {execution_time:.4f} seconds")
        
        # Return response with result and suggestions keys for frontend compatibility
        response_data = {
            "result": results,
            "suggestions": suggestions
        }
        
        return Response(response_data)
