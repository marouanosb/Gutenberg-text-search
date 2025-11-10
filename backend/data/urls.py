from data import views
from django.urls import path
from backend.config import *
from data.config import *


def construct_url_data(url):
    return URL_BASE_DATA + url

urlpatterns = [
    path('server/books/', views.BooksList.as_view()),
]