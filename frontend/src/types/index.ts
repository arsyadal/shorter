export interface URL {
  id: number;
  original_url: string;
  short_code: string;
  short_url: string;
  title: string;
  click_count: number;
  created_at: string;
}

export interface CreateURLRequest {
  url: string;
  custom_code?: string;
}

export interface CreateURLResponse {
  id: number;
  original_url: string;
  short_code: string;
  short_url: string;
  title: string;
  click_count: number;
  created_at: string;
}

export interface URLListResponse {
  urls: URL[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface ClickStats {
  total_clicks: number;
  daily_clicks: DailyClickStat[];
  country_clicks: CountryClickStat[];
  referer_clicks: RefererClickStat[];
}

export interface DailyClickStat {
  date: string;
  count: number;
}

export interface CountryClickStat {
  country: string;
  count: number;
}

export interface RefererClickStat {
  referer: string;
  count: number;
}

export interface ApiError {
  error: string;
}

export interface LoadingState {
  isLoading: boolean;
  error: string | null;
} 