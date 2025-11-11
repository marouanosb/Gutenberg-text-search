import React, { useEffect, useRef } from 'react';
import * as THREE from 'three';

// BackgroundBook: full-screen fixed Three.js scene with a stack of pages
export default function BackgroundBook({ pages = 12 }) {
  const mountRef = useRef(null);
  const rafRef = useRef(null);

  useEffect(() => {
    const mount = mountRef.current;
    if (!mount) return;

    // Renderer
    const renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true });
    renderer.setPixelRatio(window.devicePixelRatio || 1);
    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.domElement.style.position = 'fixed';
    renderer.domElement.style.left = '0';
    renderer.domElement.style.top = '0';
    renderer.domElement.style.zIndex = '0';
    renderer.domElement.style.pointerEvents = 'none';
    mount.appendChild(renderer.domElement);

    // Scene & camera
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 0.1, 1000);
    camera.position.set(0, 0, 6);

    // Lighting
    const dir = new THREE.DirectionalLight(0xffffff, 0.9);
    dir.position.set(5, 10, 7);
    scene.add(dir);
    scene.add(new THREE.AmbientLight(0xffffff, 0.6));

    // Book group
    const group = new THREE.Group();

    // Create pages as thin boxes to give some volume
    const pageGeo = new THREE.BoxGeometry(2.6, 3.6, 0.02);
    for (let i = 0; i < pages; i++) {
      const material = new THREE.MeshPhongMaterial({ color: 0xffffff, shininess: 5 });
      const page = new THREE.Mesh(pageGeo, material);
      page.position.z = -i * 0.01; // stack depth
      page.position.x = 0.0;
      page.rotation.y = 0;
      page.userData.index = i;
      group.add(page);
    }

    // Cover (slightly thicker)
    const coverGeo = new THREE.BoxGeometry(2.7, 3.8, 0.06);
    const coverMat = new THREE.MeshPhongMaterial({ color: 0x1f2937 });
    const cover = new THREE.Mesh(coverGeo, coverMat);
    cover.position.z = 0.06;
    group.add(cover);

    // Add subtle base plane behind book for depth
    const baseGeo = new THREE.PlaneGeometry(8, 6);
    const baseMat = new THREE.MeshBasicMaterial({ color: 0x071129, transparent: true, opacity: 0.6 });
    const base = new THREE.Mesh(baseGeo, baseMat);
    base.position.z = -1.2;
    base.position.y = -0.2;
    scene.add(base);

    scene.add(group);

    // Position the group slightly to the left and center vertically
    group.position.set(-1.6, 0, 0);
    group.rotation.x = -0.05;

    // Resize handler
    function onResize() {
      renderer.setSize(window.innerWidth, window.innerHeight);
      camera.aspect = window.innerWidth / window.innerHeight;
      camera.updateProjectionMatrix();
    }

    window.addEventListener('resize', onResize);

    // Scroll-driven page flip
    let lastScroll = window.scrollY || window.pageYOffset;
    const maxScroll = () => Math.max(1, document.body.scrollHeight - window.innerHeight);

    function onScroll() {
      const scrollY = window.scrollY || window.pageYOffset;
      const delta = scrollY - lastScroll;
      lastScroll = scrollY;
      const progress = Math.max(0, Math.min(1, scrollY / maxScroll()));

      // Map progress to a number of pages flipped (0..pages)
      const flipCount = progress * pages;

      // For each page, compute rotation around left edge to simulate page flip
      group.children.forEach((child) => {
        if (!child.userData.index && child.userData.index !== 0) return; // ignore cover & others
        const i = child.userData.index;
        const pageProgress = Math.max(0, Math.min(1, flipCount - i));
        // rotate up to -170deg (almost flat)
        const rot = -pageProgress * (Math.PI * 0.85);
        // ease out
        const eased = Math.sin(pageProgress * Math.PI * 0.5);
        child.rotation.y = rot * eased;
        // slight lift for pages in motion
        child.position.z = -i * 0.01 + pageProgress * 0.02;
      });

      // subtle group motion
      group.rotation.y = Math.sin(progress * Math.PI * 2) * 0.06;
      group.position.y = Math.sin(progress * Math.PI) * 0.08;
    }

    window.addEventListener('scroll', onScroll, { passive: true });

    // Render loop
    function animate() {
      rafRef.current = requestAnimationFrame(animate);
      renderer.render(scene, camera);
    }
    animate();

    // Cleanup
    return () => {
      window.removeEventListener('resize', onResize);
      window.removeEventListener('scroll', onScroll);
      cancelAnimationFrame(rafRef.current);
      renderer.dispose();
      // remove canvas
      if (renderer.domElement && mount.contains(renderer.domElement)) {
        mount.removeChild(renderer.domElement);
      }
    };
  }, [pages]);

  return <div ref={mountRef} className="background-book" style={{ position: 'fixed', inset: 0, zIndex: 0, pointerEvents: 'none' }} />;
}
