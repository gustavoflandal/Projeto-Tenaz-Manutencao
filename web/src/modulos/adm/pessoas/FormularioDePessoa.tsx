// Formulário da pessoa (RF-ADM-015): quem executa, conduz ou responde por
// algo, com ou sem usuário. O aviso de dado de saúde aparece em toda carga,
// sem botão de fechar. A identificação (CPF, CUIL, DNI ou CI) é opcional, válida
// e única na empresa; o tipo padrão segue o país da empresa e a pessoa aceita
// os de qualquer país (ADR-013).

import { ehTipoDePessoa, mascararEntrada, tiposDePessoaParaOPais, TIPOS_DE_PESSOA, type TipoDePessoa } from '@tenaz/dominio';
import { useEffect, useState, type ChangeEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { AvisoDadoDeSaude } from '../../../componentes/AvisoDadoDeSaude';
import { Botao } from '../../../componentes/Botao';
import { CampoSelecao, CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import { alteracoesNaoSalvas, tratarErroDeEnvio, type AlteracaoNaoSalva } from '../../../lib/erros';
import { dataParaISO, mascararData } from '../../../lib/formatacao';
import {
  conferirIdentificacao,
  emailOpcional,
  textoObrigatorio,
  useFormulario,
  vazioParaIndefinido,
  vazioParaNulo,
  z,
} from '../../../lib/formularios';
import type { EntradaPessoa, Pessoa } from '../../../lib/tipos';
import { VINCULOS } from '../constantes';
import { aplicarMascara, modoDeEntrada, paraOCampo, paraOEnvio } from '../identificacao';
import { useAlterarPessoa, useCriarPessoa } from '../pessoas';
import { useFormatos } from '../formatos';
import { useRegiao } from '../regiao';
import { useUnidades } from '../unidades';

const esquema = z
  .object({
  nome: textoObrigatorio('validacao.nomeObrigatorio').max(200, 'validacao.textoLongo'),
  identificacao_tipo: z.enum(TIPOS_DE_PESSOA),
  identificacao: z.string().trim(),
  data_nascimento: z
    .string()
    .trim()
    .refine((valor) => valor === '' || dataParaISO(valor) !== null, 'validacao.dataInvalida'),
  matricula: z.string().trim().max(40, 'validacao.textoLongo'),
  cargo: z.string().trim().max(100, 'validacao.textoLongo'),
  telefone: z.string().trim().max(30, 'validacao.textoLongo'),
  email: emailOpcional,
  vinculo: z.enum(VINCULOS),
  unidade_id: z.string(),
  registro_profissional: z.string().trim().max(60, 'validacao.textoLongo'),
  })
  .superRefine((v, ctx) => conferirIdentificacao(ctx, v.identificacao_tipo, v.identificacao, { obrigatoria: false }));

/** Tipo gravado na pessoa, ou o padrão do país da empresa. */
function tipoDaPessoa(pessoa: Pessoa | undefined, padrao: TipoDePessoa): TipoDePessoa {
  return ehTipoDePessoa(pessoa?.identificacao_tipo) ? pessoa.identificacao_tipo : padrao;
}

function valoresDaPessoa(pessoa: Pessoa | undefined, padrao: TipoDePessoa, formatarData: (iso: string) => string) {
  const tipo = tipoDaPessoa(pessoa, padrao);
  return {
    nome: pessoa?.nome ?? '',
    identificacao_tipo: tipo,
    identificacao: paraOCampo(tipo, pessoa?.identificacao),
    data_nascimento: pessoa?.data_nascimento ? formatarData(pessoa.data_nascimento) : '',
    matricula: pessoa?.matricula ?? '',
    cargo: pessoa?.cargo ?? '',
    telefone: pessoa?.telefone ?? '',
    email: pessoa?.email ?? '',
    vinculo: pessoa?.vinculo ?? ('empregado' as const),
    unidade_id: pessoa?.unidade_id ?? '',
    registro_profissional: pessoa?.registro_profissional ?? '',
  };
}

export interface FormularioDePessoaProps {
  /** Ausente: cria. Presente: edita. */
  readonly pessoa?: Pessoa;
  readonly podeEditar?: boolean;
  readonly aoConflitoDeVersao?: (alteracoes: readonly AlteracaoNaoSalva[]) => void;
}

export function FormularioDePessoa({ pessoa, podeEditar = true, aoConflitoDeVersao }: FormularioDePessoaProps) {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const criando = pessoa === undefined;
  const criar = useCriarPessoa();
  const alterar = useAlterarPessoa();
  const unidades = useUnidades();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [salva, setSalva] = useState(false);
  const regiao = useRegiao();
  const formatos = useFormatos();
  const tipos = tiposDePessoaParaOPais(regiao.pais);
  const padrao = tipos[0] ?? 'cpf';
  // Na criação só `defaultValues`: `values` zeraria o que a pessoa já digitou quando a região da empresa terminar de carregar.
  const formulario = useFormulario(esquema, {
    defaultValues: valoresDaPessoa(undefined, padrao, formatos.data),
    ...(pessoa ? { values: valoresDaPessoa(pessoa, padrao, formatos.data) } : {}),
  });
  const { errors, dirtyFields } = formulario.formState;
  const tipo = formulario.watch('identificacao_tipo');
  const { setValue } = formulario;
  const tipoFoiEscolhido = Boolean(dirtyFields.identificacao_tipo);
  useEffect(() => {
    // O tipo padrão segue o país da empresa assim que ele é conhecido, se a pessoa ainda não escolheu.
    if (criando && regiao.pronta && !tipoFoiEscolhido) setValue('identificacao_tipo', padrao);
  }, [criando, regiao.pronta, tipoFoiEscolhido, padrao, setValue]);
  const pendente = criar.isPending || alterar.isPending;
  const somenteLeitura = !podeEditar;

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    setSalva(false);
    const tratar = (erro: unknown) =>
      tratarErroDeEnvio(erro, {
        setError: formulario.setError,
        definirMensagem: setMensagem,
        ...(aoConflitoDeVersao ? { aoConflitoDeVersao } : {}),
        alteracoes: () =>
          alteracoesNaoSalvas(v, dirtyFields, {
            nome: t('pessoa.nome'),
            identificacao: t(`identificacao.tipos.${tipo}`),
            identificacao_tipo: t('identificacao.tipoDoDocumento'),
            data_nascimento: t('pessoa.dataNascimento'),
            vinculo: t('pessoa.vinculo'),
            unidade_id: t('pessoa.unidade'),
            matricula: t('pessoa.matricula'),
            cargo: t('pessoa.cargo'),
            telefone: t('pessoa.telefone'),
            email: t('pessoa.email'),
            registro_profissional: t('pessoa.registroProfissional'),
          }),
      });
    const nascimento = vazioParaIndefinido(v.data_nascimento) ? dataParaISO(v.data_nascimento) : null;

    if (criando) {
      const corpo: EntradaPessoa = {
        nome: v.nome,
        vinculo: v.vinculo,
        ...(vazioParaIndefinido(v.identificacao)
          ? { identificacao_tipo: v.identificacao_tipo, identificacao: paraOEnvio(v.identificacao) }
          : {}),
        ...(nascimento ? { data_nascimento: nascimento } : {}),
        ...(vazioParaIndefinido(v.matricula) ? { matricula: v.matricula } : {}),
        ...(vazioParaIndefinido(v.cargo) ? { cargo: v.cargo } : {}),
        ...(vazioParaIndefinido(v.telefone) ? { telefone: v.telefone } : {}),
        ...(vazioParaIndefinido(v.email) ? { email: v.email } : {}),
        ...(v.unidade_id ? { unidade_id: v.unidade_id } : {}),
        ...(vazioParaIndefinido(v.registro_profissional) ? { registro_profissional: v.registro_profissional } : {}),
      };
      criar.mutate(corpo, { onSuccess: (nova) => void navegar(`/administracao/pessoas/${nova.id}`), onError: tratar });
      return;
    }

    alterar.mutate(
      {
        id: pessoa.id,
        versao: pessoa.versao,
        corpo: {
          ...(dirtyFields.nome ? { nome: v.nome } : {}),
          ...(dirtyFields.identificacao || dirtyFields.identificacao_tipo
            ? vazioParaIndefinido(v.identificacao)
              ? { identificacao_tipo: v.identificacao_tipo, identificacao: paraOEnvio(v.identificacao) }
              : { identificacao_tipo: null, identificacao: null }
            : {}),
          ...(dirtyFields.data_nascimento ? { data_nascimento: nascimento } : {}),
          ...(dirtyFields.matricula ? { matricula: vazioParaNulo(v.matricula) } : {}),
          ...(dirtyFields.cargo ? { cargo: vazioParaNulo(v.cargo) } : {}),
          ...(dirtyFields.telefone ? { telefone: vazioParaNulo(v.telefone) } : {}),
          ...(dirtyFields.email ? { email: vazioParaNulo(v.email) } : {}),
          ...(dirtyFields.vinculo ? { vinculo: v.vinculo } : {}),
          ...(dirtyFields.unidade_id ? { unidade_id: v.unidade_id || null } : {}),
          ...(dirtyFields.registro_profissional ? { registro_profissional: vazioParaNulo(v.registro_profissional) } : {}),
        },
      },
      { onSuccess: () => setSalva(true), onError: tratar },
    );
  });


  // O resumo de erros diz de qual campo é cada mensagem.
  const rotulosDoResumo: Record<string, string> = {
    nome: t('pessoa.nome'),
    identificacao_tipo: t('identificacao.tipoDoDocumento'),
    identificacao: t(`identificacao.tipos.${tipo}`),
    data_nascimento: t('pessoa.dataNascimento'),
    matricula: t('pessoa.matricula'),
    cargo: t('pessoa.cargo'),
    telefone: t('pessoa.telefone'),
    email: t('pessoa.email'),
    vinculo: t('pessoa.vinculo'),
    unidade_id: t('pessoa.unidade'),
    registro_profissional: t('pessoa.registroProfissional'),
  };

  return (
    <Formulario erros={errors} rotulos={rotulosDoResumo} onSubmit={(e) => void enviar(e)} className="flex max-w-180 flex-col gap-6">
      <AvisoDadoDeSaude />
      {mensagem && <Alerta tipo="erro" rolar>{mensagem}</Alerta>}
      {salva && <Alerta tipo="sucesso" rolar>{t('pessoa.salva')}</Alerta>}
      <div className="grid gap-4 md:grid-cols-2">
        <div className="md:col-span-2">
          <CampoTexto rotulo={t('pessoa.nome')} autoComplete="off" readOnly={somenteLeitura} erro={errors.nome} {...formulario.register('nome')} />
        </div>
        <CampoSelecao
          rotulo={t('identificacao.tipoDoDocumento')}
          disabled={somenteLeitura}
          erro={errors.identificacao_tipo}
          {...formulario.register('identificacao_tipo', {
            // O mesmo número, reescrito na máscara do tipo novo.
            onChange: (e: ChangeEvent<HTMLSelectElement>) =>
              formulario.setValue('identificacao', mascararEntrada(e.target.value, formulario.getValues('identificacao'))),
          })}
        >
          {tipos.map((codigo) => (
            <option key={codigo} value={codigo}>
              {t(`identificacao.tipos.${codigo}`)}
            </option>
          ))}
        </CampoSelecao>
        <CampoTexto
          rotulo={t(`identificacao.tipos.${tipo}`)}
          apoio={t('identificacao.exemplo', { exemplo: t(`identificacao.exemplos.${tipo}`) })}
          opcional
          dado
          inputMode={modoDeEntrada(tipo)}
          autoComplete="off"
          readOnly={somenteLeitura}
          erro={errors.identificacao}
          {...formulario.register('identificacao', { onChange: aplicarMascara(() => formulario.getValues('identificacao_tipo')) })}
        />
        <CampoTexto
          rotulo={t('pessoa.dataNascimento')}
          opcional
          dado
          inputMode="numeric"
          placeholder={t('feriados.dataExemplo')}
          autoComplete="off"
          readOnly={somenteLeitura}
          erro={errors.data_nascimento}
          {...formulario.register('data_nascimento', { onChange: (e) => (e.target.value = mascararData(e.target.value)) })}
        />
        <CampoSelecao rotulo={t('pessoa.vinculo')} disabled={somenteLeitura} erro={errors.vinculo} {...formulario.register('vinculo')}>
          {VINCULOS.map((vinculo) => (
            <option key={vinculo} value={vinculo}>
              {t(`vinculos.${vinculo}`)}
            </option>
          ))}
        </CampoSelecao>
        <CampoSelecao rotulo={t('pessoa.unidade')} opcional disabled={somenteLeitura} erro={errors.unidade_id} {...formulario.register('unidade_id')}>
          <option value="">{t('pessoa.semUnidade')}</option>
          {(unidades.data?.dados ?? []).map((u) => (
            <option key={u.id} value={u.id}>
              {u.nome}
            </option>
          ))}
        </CampoSelecao>
        <CampoTexto rotulo={t('pessoa.matricula')} opcional dado readOnly={somenteLeitura} erro={errors.matricula} {...formulario.register('matricula')} />
        <CampoTexto rotulo={t('pessoa.cargo')} opcional readOnly={somenteLeitura} erro={errors.cargo} {...formulario.register('cargo')} />
        <CampoTexto rotulo={t('pessoa.telefone')} opcional type="tel" autoComplete="off" readOnly={somenteLeitura} erro={errors.telefone} {...formulario.register('telefone')} />
        <CampoTexto rotulo={t('pessoa.email')} opcional type="email" autoComplete="off" readOnly={somenteLeitura} erro={errors.email} {...formulario.register('email')} />
        <CampoTexto rotulo={t('pessoa.registroProfissional')} opcional dado readOnly={somenteLeitura} erro={errors.registro_profissional} {...formulario.register('registro_profissional')} />
      </div>
      {pessoa?.fornecedor_id && <p className="texto-corpo-p text-texto-secundario">{t('pessoa.deFornecedor')}</p>}
      {podeEditar && (
        <div className="flex justify-end gap-2">
          {criando && <Botao onClick={() => void navegar('/administracao/pessoas')}>{t('comum.voltar')}</Botao>}
          <Botao type="submit" variante="primario" carregando={pendente}>
            {pendente ? t('comum.salvando') : criando ? t('pessoa.criar') : t('comum.salvar')}
          </Botao>
        </div>
      )}
    </Formulario>
  );
}
