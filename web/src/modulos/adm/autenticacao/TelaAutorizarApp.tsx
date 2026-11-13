// Autorização do app de campo (/autorizar-app), OAuth 2.1 com PKCE
// (documento 3, seção 4.3). O app abre o navegador em GET /oauth/authorize, a
// API redireciona para cá com os parâmetros, a pessoa entra com a sessão web
// normal (senha e MFA), escolhe a empresa e esta tela chama POST /oauth/authorize.
// A resposta traz a URI de retorno com o código de uso único; a web navega até
// ela e o sistema operacional entrega ao app. O app nunca vê a senha.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { mensagemDoErro } from '../../../lib/erros';
import { EscolhaDeEmpresa } from './EscolhaDeEmpresa';
import { useAutorizarApp, useSessao } from '../sessao';

import { CartaoPublico } from './CartaoPublico';

interface PedidoDoApp {
  readonly client_id: string;
  readonly redirect_uri: string;
  readonly code_challenge: string;
  readonly state: string | undefined;
}

/** Lê os parâmetros do pedido; devolve null se faltar algum obrigatório ou se o método não for S256. */
function lerPedido(parametros: URLSearchParams): PedidoDoApp | null {
  const clientId = parametros.get('client_id');
  const redirectUri = parametros.get('redirect_uri');
  const desafio = parametros.get('code_challenge');
  if (parametros.get('response_type') !== 'code') return null;
  if (parametros.get('code_challenge_method') !== 'S256') return null;
  if (!clientId || !redirectUri || !desafio) return null;
  return { client_id: clientId, redirect_uri: redirectUri, code_challenge: desafio, state: parametros.get('state') ?? undefined };
}

export function TelaAutorizarApp() {
  const { t } = useTranslation();
  const [parametros] = useSearchParams();
  const pedido = lerPedido(parametros);
  const { data: sessao } = useSessao();
  const autorizar = useAutorizarApp();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [entregue, setEntregue] = useState(false);
  const [negado, setNegado] = useState(false);

  if (!pedido) {
    return (
      <CartaoPublico titulo={t('autorizarApp.titulo')}>
        <Alerta tipo="erro">{t('autorizarApp.pedidoInvalido')}</Alerta>
      </CartaoPublico>
    );
  }
  if (!sessao) return null;

  const enviar = (empresaId: string) => {
    setMensagem(null);
    autorizar.mutate(
      {
        response_type: 'code',
        client_id: pedido.client_id,
        redirect_uri: pedido.redirect_uri,
        code_challenge: pedido.code_challenge,
        code_challenge_method: 'S256',
        empresa_escolhida_id: empresaId,
        ...(pedido.state === undefined ? {} : { state: pedido.state }),
      },
      {
        onSuccess: ({ uri_redirecionamento }) => {
          // A URI de retorno é a que o app registrou; qualquer outra é recusada.
          if (!uri_redirecionamento.startsWith(pedido.redirect_uri)) {
            setMensagem(t('autorizarApp.destinoInesperado'));
            return;
          }
          setEntregue(true);
          globalThis.location.assign(uri_redirecionamento);
        },
        onError: (erro) => setMensagem(mensagemDoErro(erro)),
      },
    );
  };

  // "Não autorizar" devolve ao app o erro access_denied do OAuth pela URI de retorno
  // registrada (RFC 6749, 4.1.2.1); não leva ao painel da web.
  const negar = () => {
    const destino = new URL(pedido.redirect_uri);
    destino.searchParams.set('error', 'access_denied');
    if (pedido.state !== undefined) destino.searchParams.set('state', pedido.state);
    setEntregue(true);
    setNegado(true);
    globalThis.location.assign(destino.toString());
  };

  const empresaAtual = sessao.empresas?.find((e) => e.atual) ?? sessao.empresas?.[0];
  const maisDeUma = (sessao.empresas?.length ?? 0) > 1;

  return (
    <CartaoPublico titulo={t('autorizarApp.titulo')}>
      <p className="texto-corpo text-texto-secundario">{t('autorizarApp.texto', { nome: sessao.usuario.nome })}</p>
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      {entregue && <Alerta tipo="sucesso">{negado ? t('autorizarApp.negado') : t('autorizarApp.entregue')}</Alerta>}
      {maisDeUma ? (
        <EscolhaDeEmpresa
          sessao={sessao}
          aoConcluir={() => undefined}
          aoEscolher={enviar}
          rotuloConfirmar={t('autorizarApp.autorizar')}
          carregando={autorizar.isPending}
        />
      ) : (
        <Botao
          variante="primario"
          carregando={autorizar.isPending}
          disabled={!empresaAtual}
          onClick={() => empresaAtual && enviar(empresaAtual.id)}
          className="w-full"
        >
          {autorizar.isPending ? t('autorizarApp.autorizando') : t('autorizarApp.autorizar')}
        </Botao>
      )}
      <Botao variante="terciario" onClick={negar} disabled={autorizar.isPending || entregue}>
        {t('autorizarApp.negar')}
      </Botao>
    </CartaoPublico>
  );
}
