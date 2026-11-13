import { describe, expect, it } from 'vitest';

import { formatarDataHora } from './formatacao';

describe('formatarDataHora', () => {
  it('usa dd/mm/aaaa hh:mm em 24 h', () => {
    expect(formatarDataHora('2026-11-12T14:05:00Z', 'UTC')).toBe('12/11/2026 14:05');
    expect(formatarDataHora('2026-11-12T03:07:00Z', 'America/Sao_Paulo')).toBe('12/11/2026 00:07');
  });

  it('devolve null para valor vazio ou inválido', () => {
    expect(formatarDataHora('')).toBeNull();
    expect(formatarDataHora('ontem')).toBeNull();
  });
});
