// Empresa (RF-ADM-001): dados cadastrais, identificação (CNPJ, CUIT ou RUC, a do
// país) editável pelo `administrador` (a troca vai para a trilha), país e moeda
// só de leitura, idioma padrão e o interruptor do MFA obrigatório
// (RF-ADM-006, RN-ADM-027).

import { ehTipoDeEmpresa, PERFIS_DE_PAIS, type TipoDeEmpresa } from '@tenaz/dominio';
import { IDIOMAS } from '@tenaz/i18n';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Alerta } from '../../../componentes/Alerta';
import { AlertaDeConflito } from '../../../componentes/AlertaDeConflito';
import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { CampoSelecao, CampoTexto, Interruptor } from '../../../componentes/Campo';
import { ChipSituacao } from '../../../componentes/ChipStatus';
import { Confirmacao } from '../../../componentes/Dialogos';
import { Carregando, EstadoErro } from '../../../componentes/Estados';
import { Formulario } from '../../../componentes/Formulario';
import { traduzirOuOriginal } from '../../../lib/i18n';
import { alteracoesNaoSalvas, ehConflitoDeVersao, mensagemDoErro, tratarErroDeEnvio, type AlteracaoNaoSalva } from '../../../lib/erros';
import { conferirIdentificacao, textoObrigatorio, useFormulario, vazioParaNulo, z } from '../../../lib/formularios';
import type { Empresa } from '../../../lib/tipos';
import { fusosParaSeletor } from '../constantes';
import { useAlterarEmpresa, useEmpresa } from '../empresa';
import { aplicarMascara, modoDeEntrada, paraOCampo, paraOEnvio } from '../identificacao';
import { useTemPermissao } from '../sessao';

/** Tipo de identificação da empresa: o gravado, ou o do país. */
function tipoDaEmpresa(empresa: Empresa): TipoDeEmpresa {
  return ehTipoDeEmpresa(empresa.identificacao_tipo) ? empresa.identificacao_tipo : PERFIS_DE_PAIS[empresa.pais].tipo_empresa;
}

function esquemaDaEmpresa(tipo: TipoDeEmpresa) {
  return z
    .object({
      razao_social: textoObrigatorio('validacao.razaoSocialObrigatoria').max(200, 'validacao.textoLongo'),
      nome_fantasia: z.string().trim().max(200, 'validacao.textoLongo'),
      identificacao: z.string().trim(),
      idioma_padrao: z.enum(IDIOMAS),
      fuso_horario: z.string().min(1, 'validacao.obrigatorio'),
    })
    .superRefine((v, ctx) => conferirIdentificacao(ctx, tipo, v.identificacao, { obrigatoria: true }));
}

function valoresIniciais(empresa: Empresa) {
  return {
    razao_social: empresa.razao_social,
    nome_fantasia: empresa.nome_fantasia ?? '',
    identificacao: paraOCampo(tipoDaEmpresa(empresa), empresa.identificacao),
    idioma_padrao: empresa.idioma_padrao,
    fuso_horario: empresa.fuso_horario,
  };
}

function DadosCadastrais({ empresa, podeEditar }: { readonly empresa: Empresa; readonly podeEditar: boolean }) {
  const { t } = useTranslation();
  const alterar = useAlterarEmpresa();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [salvo, setSalvo] = useState(false);
  const [conflito, setConflito] = useState<readonly AlteracaoNaoSalva[] | null>(null);
  const tipo = tipoDaEmpresa(empresa);
  const formulario = useFormulario(esquemaDaEmpresa(tipo), { values: valoresIniciais(empresa) });
  const { errors, dirtyFields } = formulario.formState;

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    setSalvo(false);
    setConflito(null);
    // Só o que mudou: o PATCH não mexe no resto.
    const corpo = {
      ...(dirtyFields.razao_social ? { razao_social: v.razao_social } : {}),
      ...(dirtyFields.nome_fantasia ? { nome_fantasia: vazioParaNulo(v.nome_fantasia) } : {}),
      ...(dirtyFields.identificacao ? { identificacao_tipo: tipo, identificacao: paraOEnvio(v.identificacao) } : {}),
      ...(dirtyFields.idioma_padrao ? { idioma_padrao: v.idioma_padrao } : {}),
      ...(dirtyFields.fuso_horario ? { fuso_horario: v.fuso_horario } : {}),
    };
    alterar.mutate(
      { versao: empresa.versao, corpo },
      {
        onSuccess: () => setSalvo(true),
        onError: (erro) =>
          tratarErroDeEnvio(erro, {
            setError: formulario.setError,
            definirMensagem: setMensagem,
            aoConflitoDeVersao: setConflito,
            alteracoes: () =>
              alteracoesNaoSalvas(v, dirtyFields, {
                razao_social: t('empresa.razaoSocial'),
                nome_fantasia: t('empresa.nomeFantasia'),
                identificacao: t(`identificacao.tipos.${tipo}`),
                idioma_padrao: t('empresa.idiomaPadrao'),
                fuso_horario: t('empresa.fusoHorario'),
              }),
          }),
      },
    );
  });


  // O resumo de erros diz de qual campo é cada mensagem.
  const rotulosDoResumo: Record<string, string> = {
    razao_social: t('empresa.razaoSocial'),
    nome_fantasia: t('empresa.nomeFantasia'),
    identificacao: t(`identificacao.tipos.${tipo}`),
    idioma_padrao: t('empresa.idiomaPadrao'),
    fuso_horario: t('empresa.fusoHorario'),
  };

  return (
    <Formulario erros={errors} rotulos={rotulosDoResumo} onSubmit={(e) => void enviar(e)} className="flex max-w-140 flex-col gap-4">
      {mensagem && <Alerta tipo="erro" rolar>{mensagem}</Alerta>}
      {conflito && <AlertaDeConflito alteracoes={conflito} />}
      {salvo && <Alerta tipo="sucesso" rolar>{t('empresa.salva')}</Alerta>}
      {!podeEditar && <Alerta tipo="informacao">{t('empresa.somenteLeitura')}</Alerta>}
      <CampoTexto rotulo={t('empresa.razaoSocial')} readOnly={!podeEditar} erro={errors.razao_social} {...formulario.register('razao_social')} />
      <CampoTexto rotulo={t('empresa.nomeFantasia')} opcional readOnly={!podeEditar} erro={errors.nome_fantasia} {...formulario.register('nome_fantasia')} />
      <CampoTexto rotulo={t('pais.rotulo')} value={t(`pais.nomes.${empresa.pais}`)} readOnly apoio={t('pais.apoio')} />
      <CampoTexto rotulo={t('moeda.rotulo')} value={t(`moeda.nomes.${empresa.moeda}`)} readOnly />
      <CampoTexto
        rotulo={t(`identificacao.tipos.${tipo}`)}
        apoio={podeEditar ? t('identificacao.trocaNaTrilha') : undefined}
        inputMode={modoDeEntrada(tipo)}
        autoComplete="off"
        dado
        readOnly={!podeEditar}
        erro={errors.identificacao}
        {...formulario.register('identificacao', { onChange: aplicarMascara(() => tipo) })}
      />
      <CampoSelecao rotulo={t('empresa.idiomaPadrao')} apoio={t('empresa.idiomaPadraoApoio')} disabled={!podeEditar} erro={errors.idioma_padrao} {...formulario.register('idioma_padrao')}>
        {IDIOMAS.map((idioma) => (
          <option key={idioma} value={idioma} lang={idioma}>
            {t(`idiomas.${idioma}`)}
          </option>
        ))}
      </CampoSelecao>
      <CampoSelecao rotulo={t('empresa.fusoHorario')} disabled={!podeEditar} erro={errors.fuso_horario} {...formulario.register('fuso_horario')}>
        {fusosParaSeletor(empresa.pais, empresa.fuso_horario).map((fuso) => (
          <option key={fuso} value={fuso}>
            {traduzirOuOriginal('fusos', fuso)}
          </option>
        ))}
      </CampoSelecao>
      {podeEditar && (
        <div>
          <Botao type="submit" variante="primario" carregando={alterar.isPending}>
            {alterar.isPending ? t('comum.salvando') : t('comum.salvar')}
          </Botao>
        </div>
      )}
    </Formulario>
  );
}

function SegurancaDaEmpresa({ empresa, podeEditar }: { readonly empresa: Empresa; readonly podeEditar: boolean }) {
  const { t } = useTranslation();
  const alterar = useAlterarEmpresa();
  const [confirmando, setConfirmando] = useState(false);
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [conflito, setConflito] = useState(false);

  const gravar = (valor: boolean) => {
    setMensagem(null);
    setConflito(false);
    alterar.mutate(
      { versao: empresa.versao, corpo: { mfa_obrigatorio: valor } },
      {
        onSuccess: () => setConfirmando(false),
        onError: (erro) => {
          if (ehConflitoDeVersao(erro)) setConflito(true);
          else setMensagem(mensagemDoErro(erro));
        },
      },
    );
  };

  return (
    <div className="flex max-w-140 flex-col gap-4">
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      {conflito && <AlertaDeConflito alteracoes={[]} />}
      <Interruptor
        rotulo={t('empresa.mfaRotulo')}
        apoio={t('empresa.mfaApoio')}
        marcado={empresa.mfa_obrigatorio}
        desabilitado={!podeEditar || alterar.isPending}
        aoMudar={(marcado) => (marcado ? setConfirmando(true) : gravar(false))}
      />
      <Confirmacao
        aberto={confirmando}
        aoMudar={setConfirmando}
        titulo={t('empresa.mfaConfirmarTitulo')}
        texto={t('empresa.mfaConfirmarTexto')}
        rotuloConfirmar={t('empresa.mfaConfirmar')}
        carregando={alterar.isPending}
        aoConfirmar={() => gravar(true)}
        erro={mensagem ? <Alerta tipo="erro">{mensagem}</Alerta> : undefined}
      />
    </div>
  );
}

export function TelaEmpresa() {
  const { t } = useTranslation();
  const consulta = useEmpresa();
  const podeEditar = useTemPermissao('adm.gerenciar_empresa');

  if (consulta.isPending) {
    return (
      <>
        <CabecalhoPagina titulo={t('empresa.titulo')} />
        <Carregando rotulo={t('empresa.carregando')} linhas={4} />
      </>
    );
  }
  if (consulta.isError) return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;

  const empresa = consulta.data;
  return (
    <>
      <CabecalhoPagina
        titulo={t('empresa.titulo')}
        metadados={
          <>
            <ChipSituacao situacao={empresa.situacao === 'cancelada' || empresa.situacao === 'suspensa' ? 'inativa' : 'ativa'} />
            <span>{t('empresa.situacao', { situacao: t(`empresa.situacoes.${empresa.situacao}`) })}</span>
            <span>{empresa.plano ? t('empresa.planoContratado', { plano: traduzirOuOriginal('planos', empresa.plano) }) : t('empresa.planoServicos')}</span>
          </>
        }
      />
      <div className="flex flex-col gap-10">
        <section aria-labelledby="empresa-dados" className="flex flex-col gap-4">
          <h2 id="empresa-dados" className="texto-h3">
            {t('empresa.dados')}
          </h2>
          <DadosCadastrais empresa={empresa} podeEditar={podeEditar} />
        </section>
        <section aria-labelledby="empresa-seguranca" className="flex flex-col gap-4">
          <h2 id="empresa-seguranca" className="texto-h3">
            {t('empresa.seguranca')}
          </h2>
          <SegurancaDaEmpresa empresa={empresa} podeEditar={podeEditar} />
        </section>
      </div>
    </>
  );
}
