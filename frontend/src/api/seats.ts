import type { Seat } from '@/types';
import { apiGet } from './client';

export function fetchSeats(
    screeningId: number | string,
    signal?: AbortSignal
): Promise<Seat[]> {
    return apiGet<Seat[]>(`/screenings/${screeningId}/seats`, { signal });
}
