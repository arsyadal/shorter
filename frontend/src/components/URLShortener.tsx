'use client';

import React, { useState } from 'react';
import { toast } from 'react-hot-toast';
import { Link, Copy, ExternalLink, Settings } from 'lucide-react';
import { urlApi } from '@/utils/api';
import { copyToClipboard, isValidURL, addProtocol, validateCustomCode, cn } from '@/utils/helpers';
import type { CreateURLResponse } from '@/types';

interface URLShortenerProps {
  onURLCreated?: (url: CreateURLResponse) => void;
}

export default function URLShortener({ onURLCreated }: URLShortenerProps) {
  const [url, setUrl] = useState('');
  const [customCode, setCustomCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<CreateURLResponse | null>(null);
  const [showCustomCode, setShowCustomCode] = useState(false);
  const [errors, setErrors] = useState<{ url?: string; customCode?: string }>({});

  const validateForm = () => {
    const newErrors: { url?: string; customCode?: string } = {};

    if (!url.trim()) {
      newErrors.url = 'URL is required';
    } else if (!isValidURL(addProtocol(url))) {
      newErrors.url = 'Please enter a valid URL';
    }

    if (customCode) {
      const customCodeValidation = validateCustomCode(customCode);
      if (!customCodeValidation.isValid) {
        newErrors.customCode = customCodeValidation.error;
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    setResult(null);

    try {
      const normalizedURL = addProtocol(url.trim());
      const response = await urlApi.createURL({
        url: normalizedURL,
        custom_code: customCode.trim() || undefined,
      });

      setResult(response);
      onURLCreated?.(response);
      toast.success('Short URL created successfully!');
      
      // Reset form
      setUrl('');
      setCustomCode('');
      setShowCustomCode(false);
      setErrors({});
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to create short URL');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = async (text: string) => {
    const success = await copyToClipboard(text);
    if (success) {
      toast.success('Copied to clipboard!');
    } else {
      toast.error('Failed to copy to clipboard');
    }
  };

  const handleNewURL = () => {
    setResult(null);
    setUrl('');
    setCustomCode('');
    setShowCustomCode(false);
    setErrors({});
  };

  return (
    <div className="w-full max-w-2xl mx-auto">
      <div className="card">
        <div className="card-header">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary-100 rounded-lg">
              <Link className="w-6 h-6 text-primary-600" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-gray-900">Shorten URL</h2>
              <p className="text-sm text-gray-500">
                Create a short link from a long URL
              </p>
            </div>
          </div>
        </div>

        <div className="card-body">
          {!result ? (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-2">
                  Enter your URL
                </label>
                <input
                  type="text"
                  id="url"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  placeholder="https://example.com/very-long-url"
                  className={cn(
                    'input-field',
                    errors.url && 'border-red-300 focus:border-red-500 focus:ring-red-500'
                  )}
                  disabled={isLoading}
                />
                {errors.url && (
                  <p className="mt-1 text-sm text-red-600">{errors.url}</p>
                )}
              </div>

              <div className="flex items-center justify-between">
                <button
                  type="button"
                  onClick={() => setShowCustomCode(!showCustomCode)}
                  className="flex items-center gap-2 text-sm text-primary-600 hover:text-primary-700 transition-colors"
                  disabled={isLoading}
                >
                  <Settings className="w-4 h-4" />
                  {showCustomCode ? 'Hide' : 'Show'} Custom Code
                </button>
              </div>

              {showCustomCode && (
                <div className="animate-slide-up">
                  <label htmlFor="customCode" className="block text-sm font-medium text-gray-700 mb-2">
                    Custom Code (optional)
                  </label>
                  <input
                    type="text"
                    id="customCode"
                    value={customCode}
                    onChange={(e) => setCustomCode(e.target.value)}
                    placeholder="my-custom-code"
                    className={cn(
                      'input-field',
                      errors.customCode && 'border-red-300 focus:border-red-500 focus:ring-red-500'
                    )}
                    disabled={isLoading}
                  />
                  {errors.customCode && (
                    <p className="mt-1 text-sm text-red-600">{errors.customCode}</p>
                  )}
                  <p className="mt-1 text-xs text-gray-500">
                    3-20 characters, letters, numbers, and hyphens only
                  </p>
                </div>
              )}

              <button
                type="submit"
                disabled={isLoading || !url.trim()}
                className="btn-primary w-full"
              >
                {isLoading ? (
                  <div className="flex items-center justify-center gap-2">
                    <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Creating...
                  </div>
                ) : (
                  'Shorten URL'
                )}
              </button>
            </form>
          ) : (
            <div className="space-y-6 animate-fade-in">
              <div className="text-center">
                <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Link className="w-8 h-8 text-green-600" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  URL Shortened Successfully!
                </h3>
                <p className="text-gray-600">
                  Your short URL is ready to use
                </p>
              </div>

              <div className="bg-gray-50 rounded-lg p-4 space-y-4">
                <div>
                  <label className="block text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">
                    Short URL
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      type="text"
                      value={result.short_url}
                      readOnly
                      className="flex-1 px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm font-mono"
                    />
                    <button
                      onClick={() => handleCopy(result.short_url)}
                      className="btn-secondary p-2"
                      title="Copy to clipboard"
                    >
                      <Copy className="w-4 h-4" />
                    </button>
                    <a
                      href={result.short_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="btn-secondary p-2"
                      title="Open in new tab"
                    >
                      <ExternalLink className="w-4 h-4" />
                    </a>
                  </div>
                </div>

                {/* QR Code Section */}
                <div className="border-t border-gray-200 pt-4">
                  <label className="block text-xs font-medium text-gray-500 uppercase tracking-wide mb-3">
                    📱 QR Code
                  </label>
                  <div className="flex flex-col sm:flex-row gap-4 items-start">
                    <div className="flex-shrink-0">
                      <div className="bg-white p-3 rounded-lg border border-gray-200 inline-flex items-center justify-center">
                        <img
                          src={`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'}/qr/${result.short_code}/image`}
                          alt={`QR Code for ${result.short_code}`}
                          className="w-24 h-24 sm:w-32 sm:h-32"
                          onError={(e) => {
                            const target = e.target as HTMLImageElement;
                            target.style.display = 'none';
                            const fallback = target.nextElementSibling as HTMLElement;
                            if (fallback) fallback.style.display = 'flex';
                          }}
                        />
                        <div className="hidden w-24 h-24 sm:w-32 sm:h-32 bg-gray-100 flex items-center justify-center text-gray-400 text-xs rounded">
                          QR Loading...
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex-1 space-y-2">
                      <p className="text-sm text-gray-600">
                        Scan this QR code to quickly access your short URL
                      </p>
                      <div className="flex flex-wrap gap-2">
                        <a
                          href={`${process.env.NEXT_PUBLIC_API_URL?.replace('/api', '') || 'http://localhost:8080'}/api/qr/${result.short_code}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="inline-flex items-center gap-2 px-3 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors text-sm"
                        >
                          <ExternalLink className="w-4 h-4" />
                          View QR Page
                        </a>
                        <a
                          href={`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api'}/qr/${result.short_code}/image`}
                          download={`qr-${result.short_code}.png`}
                          className="inline-flex items-center gap-2 px-3 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors text-sm"
                        >
                          <Copy className="w-4 h-4" />
                          Download QR
                        </a>
                      </div>
                    </div>
                  </div>
                </div>

                <div>
                  <label className="block text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">
                    Original URL
                  </label>
                  <p className="px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm text-gray-600 break-all">
                    {result.original_url}
                  </p>
                </div>
              </div>

              <button
                onClick={handleNewURL}
                className="btn-primary w-full"
              >
                Create Another Short URL
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
} 