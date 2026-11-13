// Calendário de trabalho da unidade (RF-ADM-002, RN-ADM-019): a mesma regra
// do servidor. O sistema não impõe calendário; os modelos só preenchem o
// formulário e são editáveis antes de gravar.

export interface Turno {
  /** 1 = segunda ... 7 = domingo. */
  readonly dia_semana: number;
  readonly hora_inicio: string;
  /** Menor que o início termina no dia seguinte; "24:00" é meia-noite. */
  readonly hora_fim: string;
  readonly nome?: string;
}

export const DIAS_DA_SEMANA = [1, 2, 3, 4, 5, 6, 7] as const;

const RE_HORA = /^([01][0-9]|2[0-3]):[0-5][0-9]$/;

export type ErroDeTurno = 'hora_inicio_invalida' | 'hora_fim_invalida' | 'fim_igual_ao_inicio';

/** Confere os horários de um turno. Devolve o primeiro problema, ou null. */
export function validarTurno(turno: Pick<Turno, 'hora_inicio' | 'hora_fim'>): ErroDeTurno | null {
  if (!RE_HORA.test(turno.hora_inicio)) return 'hora_inicio_invalida';
  if (!RE_HORA.test(turno.hora_fim) && turno.hora_fim !== '24:00') return 'hora_fim_invalida';
  if (turno.hora_fim === turno.hora_inicio) return 'fim_igual_ao_inicio';
  return null;
}

const MINUTOS_DA_SEMANA = 7 * 24 * 60;

function minutos(hhmm: string): number {
  if (hhmm === '24:00') return 24 * 60;
  const [h, m] = hhmm.split(':');
  return Number(h) * 60 + Number(m);
}

interface Intervalo {
  readonly inicio: number;
  readonly fim: number;
}

/** Põe o turno na linha da semana, em minutos; o de domingo invade a segunda. */
function intervalos(turno: Turno): readonly Intervalo[] {
  const base = (turno.dia_semana - 1) * 24 * 60;
  const ini = minutos(turno.hora_inicio);
  let fim = minutos(turno.hora_fim);
  if (fim <= ini) fim += 24 * 60;
  const a = base + ini;
  const b = base + fim;
  if (b <= MINUTOS_DA_SEMANA) return [{ inicio: a, fim: b }];
  return [
    { inicio: a, fim: MINUTOS_DA_SEMANA },
    { inicio: 0, fim: b - MINUTOS_DA_SEMANA },
  ];
}

/** Dois turnos se sobrepõem na semana. Turnos encostados não se sobrepõem. */
export function sobrepoem(x: Turno, y: Turno): boolean {
  return intervalos(x).some((a) => intervalos(y).some((b) => a.inicio < b.fim && b.inicio < a.fim));
}

/** Índices do primeiro par de turnos que se sobrepõem, ou null. Ignora turno com hora inválida. */
export function primeiraSobreposicao(turnos: readonly Turno[]): readonly [number, number] | null {
  for (let i = 0; i < turnos.length; i++) {
    for (let j = i + 1; j < turnos.length; j++) {
      const a = turnos[i];
      const b = turnos[j];
      if (!a || !b || validarTurno(a) !== null || validarTurno(b) !== null) continue;
      if (sobrepoem(a, b)) return [i, j];
    }
  }
  return null;
}

export type ModeloDeTurnos = 'comercial' | 'continuo';

/** Comercial: segunda a sexta, 08:00 às 17:00. Contínuo: segunda a domingo, 00:00 às 24:00. */
export function turnosDoModelo(modelo: ModeloDeTurnos): readonly Turno[] {
  if (modelo === 'comercial') {
    return [1, 2, 3, 4, 5].map((dia) => ({ dia_semana: dia, hora_inicio: '08:00', hora_fim: '17:00' }));
  }
  return DIAS_DA_SEMANA.map((dia) => ({ dia_semana: dia, hora_inicio: '00:00', hora_fim: '24:00' }));
}
