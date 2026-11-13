// Formulário da unidade (RF-ADM-002): código, nome, identificação da filial (o
// tipo é o do país da empresa), fuso, endereço e, só na criação, os turnos. Depois de criada, os turnos mudam por
// operações próprias (tela de detalhe), cada uma com a versão do turno.

import { ehTipoDeEmpresa, normalizarCodigoPostal, type Pais, type TipoDeEmpresa } from '@tenaz/dominio';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoSelecao, CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import { alteracoesNaoSalvas, tratarErroDeEnvio, type AlteracaoNaoSalva } from '../../../lib/erros';
import { traduzirOuOriginal } from '../../../lib/i18n';
import {
  conferirEndereco,
  conferirIdentificacao,
  textoObrigatorio,
  useFormulario,
  vazioParaIndefinido,
  vazioParaNulo,
  z,
} from '../../../lib/formularios';
import type { Unidade } from '../../../lib/tipos';
import { fusosParaSeletor } from '../constantes';
import { aplicarMascara, modoDeEntrada, paraOCampo, paraOEnvio } from '../identificacao';
import { useRegiao } from '../regiao';
import { useAlterarUnidade, useCriarUnidade } from '../unidades';

import { CampoDeSubdivisao } from './CampoDeSubdivisao';
import { EditorDeTurnos } from './EditorDeTurnos';
import { esquemaTurnosSemExigencia, esquemaTurnos, turnosParaOCorpo } from './turnos';

function esquemaDaUnidade(criando: boolean, tipo: TipoDeEmpresa, pais: Pais | null) {
  const objeto = z.object({
    codigo: textoObrigatorio('validacao.codigoObrigatorio').max(20, 'validacao.codigoLongo'),
    nome: textoObrigatorio('validacao.nomeObrigatorio').max(200, 'validacao.textoLongo'),
    identificacao: z.string().trim(),
    fuso_horario: z.string(),
    codigo_postal: z.string().trim().max(20, 'validacao.textoLongo'),
    logradouro: z.string().trim().max(200, 'validacao.textoLongo'),
    numero: z.string().trim().max(20, 'validacao.textoLongo'),
    complemento: z.string().trim().max(100, 'validacao.textoLongo'),
    bairro: z.string().trim().max(100, 'validacao.textoLongo'),
    municipio: z.string().trim().max(100, 'validacao.textoLongo'),
    subdivisao: z.string(),
    turnos: criando ? esquemaTurnos : esquemaTurnosSemExigencia,
  });
  return objeto.superRefine((v, ctx) => {
    conferirIdentificacao(ctx, tipo, v.identificacao, { obrigatoria: false });
    conferirEndereco(ctx, pais, { subdivisao: v.subdivisao, codigoPostal: v.codigo_postal });
  });
}

const CAMPOS_DO_ENDERECO = ['codigo_postal', 'logradouro', 'numero', 'complemento', 'bairro', 'municipio', 'subdivisao'] as const;

/** O código postal vai normalizado (sem máscara, maiúsculas); os demais, aparados. */
function textoDoEndereco(campo: (typeof CAMPOS_DO_ENDERECO)[number], valor: string): string | undefined {
  const aparado = vazioParaIndefinido(valor);
  return campo === 'codigo_postal' && aparado ? normalizarCodigoPostal(aparado) : aparado;
}

function valoresDaUnidade(unidade: Unidade | undefined, tipo: TipoDeEmpresa) {
  const e = unidade?.endereco;
  return {
    codigo: unidade?.codigo ?? '',
    nome: unidade?.nome ?? '',
    identificacao: paraOCampo(tipo, unidade?.identificacao),
    fuso_horario: unidade?.fuso_horario ?? '',
    codigo_postal: e?.codigo_postal ?? '',
    logradouro: e?.logradouro ?? '',
    numero: e?.numero ?? '',
    complemento: e?.complemento ?? '',
    bairro: e?.bairro ?? '',
    municipio: e?.municipio ?? '',
    subdivisao: e?.subdivisao ?? '',
    turnos: [],
  };
}

export interface FormularioDeUnidadeProps {
  /** Ausente: cria. Presente: edita. */
  readonly unidade?: Unidade;
  readonly podeEditar?: boolean;
  /** Edição: chamado ao detectar versão desatualizada, para recarregar a unidade. */
  readonly aoConflitoDeVersao?: (alteracoes: readonly AlteracaoNaoSalva[]) => void;
}

export function FormularioDeUnidade({ unidade, podeEditar = true, aoConflitoDeVersao }: FormularioDeUnidadeProps) {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const criando = unidade === undefined;
  const criar = useCriarUnidade();
  const alterar = useAlterarUnidade();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [salvo, setSalvo] = useState(false);
  const regiao = useRegiao();
  const pais: Pais | null = regiao.pais;
  // O tipo gravado na unidade vale; sem ele, o do país da empresa (CNPJ enquanto a região não carrega).
  const tipo: TipoDeEmpresa = ehTipoDeEmpresa(unidade?.identificacao_tipo)
    ? unidade.identificacao_tipo
    : (regiao.perfil?.tipo_empresa ?? 'cnpj');
  // Na criação só `defaultValues`: `values` zeraria o que já foi digitado quando a região terminar de carregar.
  const formulario = useFormulario(esquemaDaUnidade(criando, tipo, pais), {
    defaultValues: valoresDaUnidade(undefined, tipo),
    ...(unidade ? { values: valoresDaUnidade(unidade, tipo) } : {}),
  });
  const { errors, dirtyFields } = formulario.formState;
  const rotuloDoCodigoPostal = pais === null || pais === 'BR' ? t('unidade.cep') : t('unidade.codigoPostal');
  const rotuloDaSubdivisao = t(`subdivisao.rotulo.${pais ?? 'BR'}`);
  const pendente = criar.isPending || alterar.isPending;

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    setSalvo(false);
    const tratar = (erro: unknown) =>
      tratarErroDeEnvio(erro, {
        setError: formulario.setError,
        definirMensagem: setMensagem,
        ...(aoConflitoDeVersao ? { aoConflitoDeVersao } : {}),
        alteracoes: () =>
          alteracoesNaoSalvas(v, dirtyFields, {
            codigo: t('unidade.codigo'),
            nome: t('unidade.nome'),
            identificacao: t('identificacao.daFilial', { tipo: t(`identificacao.tipos.${tipo}`) }),
            fuso_horario: t('empresa.fusoHorario'),
            codigo_postal: rotuloDoCodigoPostal,
            logradouro: t('unidade.logradouro'),
            numero: t('unidade.numero'),
            complemento: t('unidade.complemento'),
            bairro: t('unidade.bairro'),
            municipio: t('unidade.municipio'),
            subdivisao: rotuloDaSubdivisao,
          }),
        opcoes: { transformar: (campo) => campo.replace(/^endereco\./, '').replaceAll(/\[(\d+)\]/g, '.$1') },
      });

    if (criando) {
      const endereco = Object.fromEntries(
        CAMPOS_DO_ENDERECO.flatMap((campo) => {
          const valor = textoDoEndereco(campo, v[campo]);
          return valor === undefined ? [] : [[campo, valor]];
        }),
      );
      criar.mutate(
        {
          codigo: v.codigo,
          nome: v.nome,
          ...(vazioParaIndefinido(v.identificacao)
            ? { identificacao_tipo: tipo, identificacao: paraOEnvio(v.identificacao) }
            : {}),
          ...(v.fuso_horario ? { fuso_horario: v.fuso_horario } : {}),
          ...(Object.keys(endereco).length > 0 ? { endereco } : {}),
          turnos: turnosParaOCorpo(v.turnos),
        },
        { onSuccess: (nova) => void navegar(`/administracao/unidades/${nova.id}`), onError: tratar },
      );
      return;
    }

    const enderecoMudou = CAMPOS_DO_ENDERECO.some((campo) => dirtyFields[campo]);
    alterar.mutate(
      {
        id: unidade.id,
        versao: unidade.versao,
        corpo: {
          ...(dirtyFields.codigo ? { codigo: v.codigo } : {}),
          ...(dirtyFields.nome ? { nome: v.nome } : {}),
          ...(dirtyFields.identificacao
            ? vazioParaIndefinido(v.identificacao)
              ? { identificacao_tipo: tipo, identificacao: paraOEnvio(v.identificacao) }
              : { identificacao_tipo: null, identificacao: null }
            : {}),
          ...(dirtyFields.fuso_horario ? { fuso_horario: vazioParaNulo(v.fuso_horario) } : {}),
          ...(enderecoMudou
            ? { endereco: Object.fromEntries(CAMPOS_DO_ENDERECO.map((campo) => [campo, textoDoEndereco(campo, v[campo]) ?? null])) }
            : {}),
        },
      },
      { onSuccess: () => setSalvo(true), onError: tratar },
    );
  });

  const somenteLeitura = !podeEditar;

  // O resumo de erros diz de qual campo é cada mensagem.
  const rotulosDoResumo: Record<string, string> = {
    codigo: t('unidade.codigo'),
    nome: t('unidade.nome'),
    identificacao: t('identificacao.daFilial', { tipo: t(`identificacao.tipos.${tipo}`) }),
    fuso_horario: t('empresa.fusoHorario'),
    codigo_postal: rotuloDoCodigoPostal,
    logradouro: t('unidade.logradouro'),
    numero: t('unidade.numero'),
    complemento: t('unidade.complemento'),
    bairro: t('unidade.bairro'),
    municipio: t('unidade.municipio'),
    subdivisao: rotuloDaSubdivisao,
  };

  return (
    <Formulario erros={errors} rotulos={rotulosDoResumo} onSubmit={(e) => void enviar(e)} className="flex max-w-180 flex-col gap-6">
      {mensagem && <Alerta tipo="erro" rolar>{mensagem}</Alerta>}
      {salvo && <Alerta tipo="sucesso" rolar>{t('unidade.salva')}</Alerta>}
      <div className="grid gap-4 md:grid-cols-2">
        <CampoTexto rotulo={t('unidade.codigo')} apoio={t('unidade.codigoApoio')} dado readOnly={somenteLeitura} erro={errors.codigo} {...formulario.register('codigo')} />
        <CampoTexto rotulo={t('unidade.nome')} readOnly={somenteLeitura} erro={errors.nome} {...formulario.register('nome')} />
        <CampoTexto
          rotulo={t('identificacao.daFilial', { tipo: t(`identificacao.tipos.${tipo}`) })}
          opcional
          dado
          inputMode={modoDeEntrada(tipo)}
          autoComplete="off"
          readOnly={somenteLeitura}
          erro={errors.identificacao}
          {...formulario.register('identificacao', { onChange: aplicarMascara(() => tipo) })}
        />
        <CampoSelecao rotulo={t('empresa.fusoHorario')} disabled={somenteLeitura} erro={errors.fuso_horario} {...formulario.register('fuso_horario')}>
          <option value="">{t('unidade.fusoDaEmpresa')}</option>
          {fusosParaSeletor(pais, unidade?.fuso_horario).map((fuso) => (
            <option key={fuso} value={fuso}>
              {traduzirOuOriginal('fusos', fuso)}
            </option>
          ))}
        </CampoSelecao>
      </div>

      <fieldset className="grid gap-4 md:grid-cols-2">
        <legend className="mb-2 texto-h4">{t('unidade.endereco')}</legend>
        <CampoTexto rotulo={rotuloDoCodigoPostal} apoio={t('identificacao.exemplo', { exemplo: t(`unidade.codigoPostalExemplos.${pais ?? 'BR'}`) })} opcional dado autoComplete="postal-code" readOnly={somenteLeitura} erro={errors.codigo_postal} {...formulario.register('codigo_postal')} />
        <CampoTexto rotulo={t('unidade.logradouro')} opcional autoComplete="address-line1" readOnly={somenteLeitura} erro={errors.logradouro} {...formulario.register('logradouro')} />
        <CampoTexto rotulo={t('unidade.numero')} opcional readOnly={somenteLeitura} erro={errors.numero} {...formulario.register('numero')} />
        <CampoTexto rotulo={t('unidade.complemento')} opcional readOnly={somenteLeitura} erro={errors.complemento} {...formulario.register('complemento')} />
        <CampoTexto rotulo={t('unidade.bairro')} opcional readOnly={somenteLeitura} erro={errors.bairro} {...formulario.register('bairro')} />
        <CampoTexto rotulo={pais === null || pais === 'BR' ? t('unidade.municipio') : t('unidade.cidade')} opcional readOnly={somenteLeitura} erro={errors.municipio} {...formulario.register('municipio')} />
        <CampoDeSubdivisao pais={pais} atual={unidade?.endereco.subdivisao} disabled={somenteLeitura} erro={errors.subdivisao} {...formulario.register('subdivisao')} />
      </fieldset>

      {criando && <EditorDeTurnos control={formulario.control} register={formulario.register} errors={errors} />}

      {podeEditar && (
        <div className="flex justify-end gap-2">
          {criando && (
            <Botao onClick={() => void navegar('/administracao/unidades')}>{t('comum.voltar')}</Botao>
          )}
          <Botao type="submit" variante="primario" carregando={pendente}>
            {pendente ? t('comum.salvando') : criando ? t('unidade.criar') : t('comum.salvar')}
          </Botao>
        </div>
      )}
    </Formulario>
  );
}
