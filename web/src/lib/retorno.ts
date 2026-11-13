// Destino depois do login. Só caminhos internos da própria aplicação: um
// endereço de outra origem (`//outro.site`, `https://...`) é descartado, para
// que a tela de login não sirva de redirecionamento aberto.

export function retornoSeguro(valor: string | null | undefined, padrao = '/'): string {
  if (!valor) return padrao;
  if (!valor.startsWith('/') || valor.startsWith('//') || valor.includes('\\')) return padrao;
  if (valor.startsWith('/entrar')) return padrao;
  return valor;
}
