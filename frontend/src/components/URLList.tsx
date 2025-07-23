'use client';

import React, { useState, useEffect } from 'react';
import { toast } from 'react-hot-toast';
import { 
  ExternalLink, 
  Copy, 
  BarChart3, 
  Calendar,
  Eye,
  Globe,
  ChevronLeft,
  ChevronRight
} from 'lucide-react';
import { urlApi } from '@/utils/api';
import { 
  copyToClipboard, 
  formatRelativeTime, 
  formatNumber, 
  truncateText, 
  getDomain,
  cn 
} from '@/utils/helpers';
import type { URL, URLListResponse } from '@/types';

interface URLListProps {
  refreshTrigger?: number;
}

export default function URLList({ refreshTrigger }: URLListProps) {
  const [data, setData] = useState<URLListResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedURL, setSelectedURL] = useState<string | null>(null);

  const fetchURLs = async (page: number = 1) => {
    try {
      setIsLoading(true);
      const response = await urlApi.getURLs(page, 10);
      setData(response);
    } catch (error) {
      toast.error('Failed to fetch URLs');
      console.error('Failed to fetch URLs:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchURLs(currentPage);
  }, [currentPage, refreshTrigger]);

  const handleCopy = async (text: string, shortCode: string) => {
    const success = await copyToClipboard(text);
    if (success) {
      setSelectedURL(shortCode);
      toast.success('Copied to clipboard!');
      setTimeout(() => setSelectedURL(null), 1000);
    } else {
      toast.error('Failed to copy to clipboard');
    }
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  if (isLoading) {
    return (
      <div className="w-full max-w-6xl mx-auto">
        <div className="card">
          <div className="card-header">
            <h3 className="text-lg font-semibold text-gray-900">Recent URLs</h3>
          </div>
          <div className="card-body">
            <div className="space-y-4">
              {[...Array(5)].map((_, i) => (
                <div key={i} className="animate-pulse">
                  <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                    <div className="flex-1 space-y-2">
                      <div className="h-4 bg-gray-300 rounded w-3/4"></div>
                      <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                    </div>
                    <div className="flex items-center gap-2">
                      <div className="h-8 w-20 bg-gray-300 rounded"></div>
                      <div className="h-8 w-8 bg-gray-300 rounded"></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!data || data.urls.length === 0) {
    return (
      <div className="w-full max-w-6xl mx-auto">
        <div className="card">
          <div className="card-header">
            <h3 className="text-lg font-semibold text-gray-900">Recent URLs</h3>
          </div>
          <div className="card-body">
            <div className="text-center py-12">
              <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <BarChart3 className="w-8 h-8 text-gray-400" />
              </div>
              <h4 className="text-lg font-medium text-gray-900 mb-2">No URLs yet</h4>
              <p className="text-gray-500">
                Create your first short URL to see it here
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full max-w-6xl mx-auto">
      <div className="card">
        <div className="card-header">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-900">Recent URLs</h3>
            <div className="text-sm text-gray-500">
              {formatNumber(data.total)} total URLs
            </div>
          </div>
        </div>

        <div className="card-body p-0">
          <div className="divide-y divide-gray-200">
            {data.urls.map((url) => (
              <div key={url.id} className="p-6 hover:bg-gray-50 transition-colors">
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-2">
                      <div className="flex-shrink-0">
                        <Globe className="w-4 h-4 text-gray-400" />
                      </div>
                      <h4 className="text-sm font-medium text-gray-900 truncate">
                        {url.title || getDomain(url.original_url)}
                      </h4>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex items-center gap-2">
                        <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                          Short:
                        </span>
                        <code className="text-sm text-primary-600 font-mono bg-primary-50 px-2 py-1 rounded">
                          {url.short_url}
                        </code>
                        <button
                          onClick={() => handleCopy(url.short_url, url.short_code)}
                          className={cn(
                            "p-1 rounded hover:bg-gray-200 transition-colors",
                            selectedURL === url.short_code && "bg-green-100 text-green-600"
                          )}
                          title="Copy to clipboard"
                        >
                          <Copy className="w-3 h-3" />
                        </button>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                          Original:
                        </span>
                        <p className="text-sm text-gray-600 truncate max-w-md" title={url.original_url}>
                          {truncateText(url.original_url, 60)}
                        </p>
                        <a
                          href={url.original_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="p-1 rounded hover:bg-gray-200 transition-colors text-gray-400 hover:text-gray-600"
                          title="Open original URL"
                        >
                          <ExternalLink className="w-3 h-3" />
                        </a>
                      </div>
                    </div>

                    <div className="flex items-center gap-4 mt-3 text-xs text-gray-500">
                      <div className="flex items-center gap-1">
                        <Calendar className="w-3 h-3" />
                        <span>{formatRelativeTime(url.created_at)}</span>
                      </div>
                      <div className="flex items-center gap-1">
                        <Eye className="w-3 h-3" />
                        <span>{formatNumber(url.click_count)} clicks</span>
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <a
                      href={url.short_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="btn-secondary text-xs"
                    >
                      Visit
                    </a>
                    <button className="btn-secondary text-xs" title="View statistics">
                      <BarChart3 className="w-3 h-3" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {data.total_pages > 1 && (
          <div className="card-header border-t border-gray-200 bg-white">
            <div className="flex items-center justify-between">
              <p className="text-sm text-gray-700">
                Showing page {data.page} of {data.total_pages} ({formatNumber(data.total)} total)
              </p>
              
              <div className="flex items-center gap-2">
                <button
                  onClick={() => handlePageChange(data.page - 1)}
                  disabled={data.page <= 1}
                  className="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <ChevronLeft className="w-4 h-4" />
                  Previous
                </button>
                
                <div className="flex items-center gap-1">
                  {[...Array(Math.min(5, data.total_pages))].map((_, i) => {
                    const pageNum = Math.max(1, data.page - 2) + i;
                    if (pageNum > data.total_pages) return null;
                    
                    return (
                      <button
                        key={pageNum}
                        onClick={() => handlePageChange(pageNum)}
                        className={cn(
                          'px-3 py-1 text-sm rounded transition-colors',
                          pageNum === data.page
                            ? 'bg-primary-600 text-white'
                            : 'hover:bg-gray-100 text-gray-700'
                        )}
                      >
                        {pageNum}
                      </button>
                    );
                  })}
                </div>
                
                <button
                  onClick={() => handlePageChange(data.page + 1)}
                  disabled={data.page >= data.total_pages}
                  className="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
} 