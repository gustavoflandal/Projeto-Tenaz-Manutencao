// Esquema do calendário de trabalho no formulário (RN-ADM-019): a mesma regra
// do servidor, vinda de @tenaz/dominio. O servidor continua decidindo; o
// formulário só avisa antes de enviar.

import { primeiraSobreposicao, validarTurno, type ErroDeTurno } from '@tenaz/dominio';

import { z } from '../../../lib/formularios';

/** Turno como digitado no formulário (horas são texto). */
export interface TurnoDoFormulario {
  dia_semana: number;
  nome: string;
  hora_inicio: string;
  hora_fim: string;
}

const MENSAGEM: Record<ErroDeTurno, { campo: 'hora_inicio' | 'hora_fim'; chave: string }> = {
  hora_inicio_invalida: { campo: 'hora_inicio', chave: 'validacao.horaInicio' },
  hora_fim_invalida: { campo: 'hora_fim', chave: 'validacao.horaFim' },
  fim_igual_ao_inicio: { campo: 'hora_fim', chave: 'validacao.fimIgualAoInicio' },
};

export const esquemaTurno = z.object({
  dia_semana: z.number().int().min(1).max(7),
  nome: z.string().trim().max(60, 'validacao.nomeDoTurnoLongo'),
  hora_inicio: z.string().trim(),
  hora_fim: z.string().trim(),
});

/**
 * Confere uma lista de turnos: horários válidos e sem sobreposição; com
 * `exigirUm`, também ao menos um turno (RN-ADM-019: toda unidade tem turno).
 */
export function verificarTurnos(
  turnos: ReadonlyArray<{ dia_semana: number; hora_inicio: string; hora_fim: string }>,
  contexto: z.RefinementCtx,
  exigirUm: boolean,
): void {
  if (exigirUm && turnos.length === 0) {
    contexto.addIssue({ code: 'custom', path: [], message: 'validacao.turnoObrigatorio' });
  }
  turnos.forEach((turno, indice) => {
    const erro = validarTurno(turno);
    if (erro) contexto.addIssue({ code: 'custom', path: [indice, MENSAGEM[erro].campo], message: MENSAGEM[erro].chave });
  });
  const par = primeiraSobreposicao(turnos);
  if (par) contexto.addIssue({ code: 'custom', path: [par[1]], message: 'validacao.turnoSobrepoe' });
}

/** Lista de turnos de uma unidade nova: ao menos um, horários válidos, sem sobreposição. */
export const esquemaTurnos = z.array(esquemaTurno).superRefine((turnos, contexto) => verificarTurnos(turnos, contexto, true));

/** Lista de turnos de um formulário que não os edita (edição da unidade). */
export const esquemaTurnosSemExigencia = z.array(esquemaTurno);

/** Corpo de EntradaTurno a partir do que foi digitado. */
export function turnosParaOCorpo(turnos: readonly TurnoDoFormulario[]) {
  return turnos.map((t) => ({
    dia_semana: t.dia_semana,
    hora_inicio: t.hora_inicio,
    hora_fim: t.hora_fim,
    ...(t.nome ? { nome: t.nome } : {}),
  }));
}

/** O modelo como linhas do formulário. */
export function turnosDoFormulario(
  turnos: ReadonlyArray<{ dia_semana: number; hora_inicio: string; hora_fim: string; nome?: string | null }>,
): TurnoDoFormulario[] {
  return turnos.map((t) => ({
    dia_semana: t.dia_semana,
    nome: t.nome ?? '',
    hora_inicio: t.hora_inicio,
    hora_fim: t.hora_fim,
  }));
}
