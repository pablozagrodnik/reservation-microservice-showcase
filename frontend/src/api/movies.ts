import type { Movie } from '@/types';
import { apiGet } from './client';

export function fetchMovies(signal?: AbortSignal): Promise<Movie[]> {
    return apiGet<Movie[]>('/movies', { signal });
}
