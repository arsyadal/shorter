import React, { useState } from 'react';
import Head from 'next/head';
import { Link2, Github, Twitter } from 'lucide-react';
import URLShortener from '@/components/URLShortener';
import URLList from '@/components/URLList';

export default function Home() {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleURLCreated = () => {
    // Trigger refresh of URL list
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <>
      <Head>
        <title>URL Shortener - Create Short Links</title>
        <meta name="description" content="Create short links from long URLs. Track clicks and analyze your link performance." />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link 
          href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" 
          rel="stylesheet" 
        />
      </Head>

      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
        {/* Header */}
        <header className="border-b border-white/20 bg-white/50 backdrop-blur-sm">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex items-center justify-between h-16">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-primary-600 rounded-lg">
                  <Link2 className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h1 className="text-xl font-bold text-gray-900">Shorter</h1>
                  <p className="text-xs text-gray-500">URL Shortener</p>
                </div>
              </div>
              
              <div className="flex items-center gap-4">
                <a
                  href="https://github.com"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="p-2 text-gray-600 hover:text-gray-900 rounded-lg hover:bg-white/50 transition-colors"
                  title="GitHub Repository"
                >
                  <Github className="w-5 h-5" />
                </a>
                <a
                  href="https://twitter.com"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="p-2 text-gray-600 hover:text-gray-900 rounded-lg hover:bg-white/50 transition-colors"
                  title="Follow us on Twitter"
                >
                  <Twitter className="w-5 h-5" />
                </a>
              </div>
            </div>
          </div>
        </header>

        {/* Hero Section */}
        <section className="py-16 px-4 sm:px-6 lg:px-8">
          <div className="max-w-4xl mx-auto text-center">
            <div className="mb-8">
              <h2 className="text-4xl font-bold text-gray-900 mb-4 sm:text-5xl lg:text-6xl">
                Shorten Your URLs
                <span className="block text-primary-600">Make Them Memorable</span>
              </h2>
              <p className="text-xl text-gray-600 max-w-2xl mx-auto leading-relaxed">
                Transform long, complex URLs into short, shareable links. 
                Track clicks, analyze performance, and customize your links.
              </p>
            </div>

            {/* Features */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
              <div className="bg-white/60 backdrop-blur-sm rounded-xl p-6 border border-white/20">
                <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Link2 className="w-6 h-6 text-blue-600" />
                </div>
                <h3 className="font-semibold text-gray-900 mb-2">Easy to Use</h3>
                <p className="text-gray-600 text-sm">
                  Paste your long URL and get a short link instantly
                </p>
              </div>
              
              <div className="bg-white/60 backdrop-blur-sm rounded-xl p-6 border border-white/20">
                <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                </div>
                <h3 className="font-semibold text-gray-900 mb-2">Analytics</h3>
                <p className="text-gray-600 text-sm">
                  Track clicks and monitor link performance
                </p>
              </div>
              
              <div className="bg-white/60 backdrop-blur-sm rounded-xl p-6 border border-white/20">
                <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                </div>
                <h3 className="font-semibold text-gray-900 mb-2">Secure</h3>
                <p className="text-gray-600 text-sm">
                  Safe and reliable URL shortening service
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Main Content */}
        <main className="pb-16 px-4 sm:px-6 lg:px-8">
          <div className="space-y-12">
            {/* URL Shortener */}
            <URLShortener onURLCreated={handleURLCreated} />
            
            {/* Recent URLs */}
            <URLList refreshTrigger={refreshTrigger} />
          </div>
        </main>

        {/* Footer */}
        <footer className="border-t border-gray-200 bg-white/50 backdrop-blur-sm">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="text-center">
              <div className="flex items-center justify-center gap-2 mb-4">
                <div className="p-1 bg-primary-600 rounded">
                  <Link2 className="w-4 h-4 text-white" />
                </div>
                <span className="font-semibold text-gray-900">Shorter</span>
              </div>
              <p className="text-gray-600 text-sm mb-4">
                A simple and fast URL shortening service
              </p>
              <div className="flex items-center justify-center gap-6 text-sm text-gray-500">
                <a href="#" className="hover:text-gray-900 transition-colors">
                  Privacy Policy
                </a>
                <a href="#" className="hover:text-gray-900 transition-colors">
                  Terms of Service
                </a>
                <a href="#" className="hover:text-gray-900 transition-colors">
                  API Documentation
                </a>
              </div>
              <p className="text-xs text-gray-400 mt-4">
                © 2024 Shorter. Built with ❤️ using Go and Next.js
              </p>
            </div>
          </div>
        </footer>
      </div>
    </>
  );
} 