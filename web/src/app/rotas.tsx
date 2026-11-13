// Rotas da aplicação principal. O link de uso único (/l/#<token>) não está
// aqui: tem entrada e pacote próprios, em src/link/ (seção 6.7 do documento 5).

import type { RouteObject } from 'react-router';

import { TelaAuditoria } from '../modulos/adm/auditoria/TelaAuditoria';
import { TelaCapturas } from '../modulos/adm/auditoria/TelaCapturas';
import { TelaIntegridade } from '../modulos/adm/auditoria/TelaIntegridade';
import { TelaAtivarMFAObrigatorio } from '../modulos/adm/autenticacao/TelaAtivarMFAObrigatorio';
import { TelaAutorizarApp } from '../modulos/adm/autenticacao/TelaAutorizarApp';
import { TelaConta } from '../modulos/adm/autenticacao/TelaConta';
import { TelaConvite } from '../modulos/adm/autenticacao/TelaConvite';
import { TelaCriarConta } from '../modulos/adm/autenticacao/TelaCriarConta';
import { TelaEntrar } from '../modulos/adm/autenticacao/TelaEntrar';
import { TelaEsqueciSenha } from '../modulos/adm/autenticacao/TelaEsqueciSenha';
import { TelaRedefinirSenha } from '../modulos/adm/autenticacao/TelaRedefinirSenha';
import { TelaEmpresa } from '../modulos/adm/empresa/TelaEmpresa';
import { TelaPerfis } from '../modulos/adm/perfis/TelaPerfis';
import { TelaPessoaDetalhe, TelaPessoaNova, TelaPessoas } from '../modulos/adm/pessoas/TelaPessoas';
import { TelaUnidadeDetalhe } from '../modulos/adm/unidades/TelaUnidadeDetalhe';
import { TelaUnidadeNova, TelaUnidades } from '../modulos/adm/unidades/TelaUnidades';
import { TelaUsuarioDetalhe } from '../modulos/adm/usuarios/TelaUsuarioDetalhe';
import { TelaUsuarioNovo, TelaUsuarios } from '../modulos/adm/usuarios/TelaUsuarios';
import { TelaPainel } from '../modulos/ind/TelaPainel';

import { LayoutPrincipal } from './LayoutPrincipal';
import { RequerPermissao } from './RequerPermissao';
import { RequerSessao } from './RequerSessao';
import { TelaNaoEncontrada } from './TelaNaoEncontrada';

const criarUnidade = (
  <RequerPermissao permissao="adm.gerenciar_empresa" />
);

export const rotas: RouteObject[] = [
  // Sem sessão
  { path: '/entrar', element: <TelaEntrar /> },
  { path: '/criar-conta', element: <TelaCriarConta /> },
  { path: '/esqueci-a-senha', element: <TelaEsqueciSenha /> },
  { path: '/redefinir-senha', element: <TelaRedefinirSenha /> },
  { path: '/convite', element: <TelaConvite /> },

  // Sessão restrita à ativação do segundo fator (RN-ADM-027)
  {
    element: <RequerSessao somenteRestrita />,
    children: [{ path: '/ativar-mfa', element: <TelaAtivarMFAObrigatorio /> }],
  },

  // Com sessão plena
  {
    element: <RequerSessao />,
    children: [
      { path: '/autorizar-app', element: <TelaAutorizarApp /> },
      {
        path: '/',
        element: <LayoutPrincipal />,
        children: [
          { index: true, element: <TelaPainel /> },
          { path: 'conta', element: <TelaConta /> },
          {
            path: 'administracao',
            children: [
              {
                element: <RequerPermissao permissao="adm.ler" />,
                children: [
                  { path: 'empresa', element: <TelaEmpresa /> },
                  { path: 'unidades', element: <TelaUnidades /> },
                  { path: 'unidades/nova', element: criarUnidade, children: [{ index: true, element: <TelaUnidadeNova /> }] },
                  { path: 'unidades/:id', element: <TelaUnidadeDetalhe /> },
                ],
              },
              {
                element: <RequerPermissao permissao="adm.ler_pessoas" />,
                children: [
                  { path: 'pessoas', element: <TelaPessoas /> },
                  {
                    path: 'pessoas/nova',
                    element: <RequerPermissao permissao="adm.gerenciar_pessoas" />,
                    children: [{ index: true, element: <TelaPessoaNova /> }],
                  },
                  { path: 'pessoas/:id', element: <TelaPessoaDetalhe /> },
                ],
              },
              {
                element: <RequerPermissao permissao="adm.ler_usuarios" />,
                children: [
                  { path: 'usuarios', element: <TelaUsuarios /> },
                  {
                    path: 'usuarios/novo',
                    element: <RequerPermissao permissao="adm.gerenciar_usuarios" />,
                    children: [{ index: true, element: <TelaUsuarioNovo /> }],
                  },
                  { path: 'usuarios/:id', element: <TelaUsuarioDetalhe /> },
                  { path: 'perfis', element: <TelaPerfis /> },
                ],
              },
              {
                element: <RequerPermissao permissao="adm.ler_auditoria" />,
                children: [
                  { path: 'auditoria', element: <TelaAuditoria /> },
                  { path: 'auditoria/integridade', element: <TelaIntegridade /> },
                ],
              },
              {
                element: <RequerPermissao permissao="adm.gerenciar_parametros" />,
                children: [{ path: 'auditoria/capturas', element: <TelaCapturas /> }],
              },
            ],
          },
          { path: '*', element: <TelaNaoEncontrada /> },
        ],
      },
    ],
  },
];
