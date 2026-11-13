// Formatação de valores para exibição, sempre por Intl com pt-BR (seção 3.11
// de docs/7_PADROES_DESIGN.md). Quando o app precisar das mesmas funções,
// elas sobem para pacotes/dominio.

const LOCALIDADE = 'pt-BR';

/**
 * Data e hora no formato dd/mm/aaaa hh:mm (24 h). Devolve null quando o valor
 * não é um instante válido.
 *
 * @param fuso fuso de exibição; sem ele, vale o do navegador.
 */
export function formatarDataHora(instante: string, fuso?: string): string | null {
  const data = new Date(instante);
  if (instante.trim() === '' || Number.isNaN(data.getTime())) return null;
  const dia = new Intl.DateTimeFormat(LOCALIDADE, {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    timeZone: fuso,
  }).format(data);
  const hora = new Intl.DateTimeFormat(LOCALIDADE, {
    hour: '2-digit',
    minute: '2-digit',
    hourCycle: 'h23',
    timeZone: fuso,
  }).format(data);
  return `${dia} ${hora}`;
}
