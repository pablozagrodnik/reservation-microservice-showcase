import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Screening, Seat } from '@/types';

export const useReservationStore = defineStore('reservation', () => {
    const screening = ref<Screening | null>(null);
    const selectedSeats = ref<Seat[]>([]);

    const TICKET_PRICE = 25.00;

    const totalAmount = computed(() => selectedSeats.value.length * TICKET_PRICE);
    const isReadyForCheckout = computed(() => screening.value !== null && selectedSeats.value.length > 0);

    const setScreening = (newScreening: Screening) => {
        screening.value = newScreening;
        selectedSeats.value = [];
    };

    const toggleSeat = (seat: Seat) => {
        if (seat.is_taken) return;

        const index = selectedSeats.value.findIndex(s => s.id === seat.id);
        if (index === -1) {
            selectedSeats.value.push(seat);
        } else {
            selectedSeats.value.splice(index, 1);
        }
    };

    const clearReservation = () => {
        screening.value = null;
        selectedSeats.value = [];
    };

    return {
        screening,
        selectedSeats,
        totalAmount,
        isReadyForCheckout,
        setScreening,
        toggleSeat,
        clearReservation
    };
});