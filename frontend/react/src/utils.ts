import Cookies from 'js-cookie'

// export function getCSRFToken(): string | null {
//     const name = 'csrftoken';
//     const cookieValue = document.cookie
//         .split('; ')
//         .find(row => row.startsWith(name))
//         ?.split('=')[1];
//     return cookieValue || null;
// }

export function getCSRFToken() {
    return Cookies.get('csrftoken');
}

/**
 * @param duration = 65
 * @returns hour = 1
 */
export function hourOfDuration(duration: number | undefined | null): number {
    if (duration == null || duration < 0) {
        return 0;
    }
    return Math.floor(duration / 60);
}

/**
 * @param duration = 65
 * @returns minute = 5
 */
export function minOfDuration(duration: number | undefined | null): number {
    if (duration == null) {
        return 0;
    } else if (duration < 0) {
        return -1;
    }
    return duration % 60
}

export function percentage(current: number, total: number): number {
    switch (total) {
        case -1:    // Endless
            return -1;
        case 0:     // Undefined
            return 0;
        default:    // Defined
            if (current >= total) {
                return 100;
            } else {
                return Math.round((current / total) * 100);
            }
    }
}

export function rating(num: string): string {
    switch (num) {
        case '10':
            return '10: Perfect'
        case '9':
            return '9: Excellent'
        case '8':
            return '8: Fantastic'
        case '7':
            return '7: Great'
        case '6':
            return '6: Good'
        case '5':
            return '5: Fine'
        case '4':
            return '4: Not Satisfied'
        case '3':
            return '3: Boring'
        case '2':
            return '2: Bad'
        case '1':
            return '1: Trash'
        default:
            return ''
    }
}