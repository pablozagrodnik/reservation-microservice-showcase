import { onMounted, onUnmounted, ref } from 'vue';
import type { Movie } from '@/types';
import { fetchMovies } from '@/api/movies';
import { ApiError } from '@/api/client';

export function useMovies() {
    const movies = ref<Movie[]>([]);
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
            const data = await fetchMovies(controller.signal);
            if (activeController !== controller) return;
            movies.value = data;
        } catch (err) {
            if (activeController !== controller) return;
            if (err instanceof DOMException && err.name === 'AbortError') return;

            error.value = err instanceof ApiError
                ? `Nie udało się pobrać repertuaru (kod ${err.status}).`
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

    onUnmounted(() => {
        activeController?.abort();
    });

    return { movies, isLoading, error, refresh: load };
}
