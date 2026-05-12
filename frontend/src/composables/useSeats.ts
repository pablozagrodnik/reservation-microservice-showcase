import { onMounted, onUnmounted, ref, watch, type Ref } from 'vue';
import type { Seat } from '@/types';
import { fetchSeats } from '@/api/seats';
import { ApiError } from '@/api/client';

export function useSeats(screeningId: Ref<string | number>) {
    const seats = ref<Seat[]>([]);
    const isLoading = ref<boolean>(false);
    const error = ref<string | null>(null);

    let activeController: AbortController | null = null;

    const load = async (): Promise<void> => {
        activeController?.abort();
        const controller = new AbortController();
        activeController = controller;

        isLoading.value = true;
        error.value = null;

        try {
            const data = await fetchSeats(screeningId.value, controller.signal);
            if (activeController !== controller) return;
            seats.value = data;
        } catch (err) {
            if (activeController !== controller) return;
            if (err instanceof DOMException && err.name === 'AbortError') return;

            error.value = err instanceof ApiError
                ? `Nie udało się pobrać miejsc (kod ${err.status}).`
                : 'Nie udało się połączyć z serwerem. Sprawdź połączenie i spróbuj ponownie.';
        } finally {
            if (activeController === controller) {
                isLoading.value = false;
            }
        }
    };

    onMounted(() => {
        void load();
    });

    watch(screeningId, () => {
        void load();
    });

    onUnmounted(() => {
        activeController?.abort();
    });

    return { seats, isLoading, error, refresh: load };
}
