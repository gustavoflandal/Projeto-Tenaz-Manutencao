// Cadastro da empresa pelo site (RF-ADM-001, POST /v1/empresas).
//
// O responsável nasce com os perfis `administrador` e `gestor` e recebe um
// convite por e-mail para definir a senha. A unidade inicial pede os turnos:
// o sistema não impõe calendário, e os modelos só preenchem o formulário.

import { normalizarCodigoPostal, PAISES, PERFIS_DE_PAIS, type Pais } from '@tenaz/dominio';
import { IDIOMAS, type Idioma } from '@tenaz/i18n';
import { useState, type ChangeEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoCaixa, CampoSelecao, CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import type { NovaEmpresa } from '../../../lib/tipos';
import { aplicarErrosDoServidor, mensagemDoErro } from '../../../lib/erros';
import { traduzirOuOriginal } from '../../../lib/i18n';
import { paisDoNavegador } from '../../../lib/idioma';
import {
  conferirEndereco,
  conferirIdentificacao,
  emailObrigatorio,
  textoObrigatorio,
  useFormulario,
  vazioParaIndefinido,
  z,
} from '../../../lib/formularios';
import { useCriarEmpresa } from '../sessao';
import { CampoDeSubdivisao } from '../unidades/CampoDeSubdivisao';
import { EditorDeTurnos } from '../unidades/EditorDeTurnos';
import { esquemaTurnos, turnosParaOCorpo } from '../unidades/turnos';

import { CartaoPublico } from './CartaoPublico';
import { fusosParaSeletor, PLANOS } from '../constantes';
import { aplicarMascara, modoDeEntrada, paraOEnvio } from '../identificacao';

const esquema = z
  .object({
    razao_social: textoObrigatorio('validacao.razaoSocialObrigatoria').max(200, 'validacao.textoLongo'),
    nome_fantasia: z.string().trim().max(200, 'validacao.textoLongo'),
    pais: z.enum(PAISES),
    identificacao: z.string().trim(),
    idioma_padrao: z.enum(IDIOMAS),
    fuso_horario: z.string(),
    plano: z.string(),
    srv_habilitado: z.boolean(),
    responsavel_nome: textoObrigatorio('validacao.nomeObrigatorio').max(200, 'validacao.textoLongo'),
    responsavel_email: emailObrigatorio,
    responsavel_telefone: z.string().trim().max(30, 'validacao.textoLongo'),
    unidade_codigo: textoObrigatorio('validacao.codigoObrigatorio').max(20, 'validacao.codigoLongo'),
    unidade_nome: textoObrigatorio('validacao.nomeObrigatorio').max(200, 'validacao.textoLongo'),
    unidade_subdivisao: z.string(),
    unidade_codigo_postal: z.string().trim().max(20, 'validacao.textoLongo'),
    turnos: esquemaTurnos,
  })
  .refine((v) => v.srv_habilitado || v.plano !== '', { path: ['plano'], message: 'validacao.planoObrigatorio' })
  .superRefine((v, ctx) => {
    conferirIdentificacao(ctx, PERFIS_DE_PAIS[v.pais].tipo_empresa, v.identificacao, { obrigatoria: true });
    conferirEndereco(ctx, v.pais, { subdivisao: v.unidade_subdivisao, codigoPostal: v.unidade_codigo_postal }, 'unidade_codigo_postal');
  });

/** Nome do campo do servidor (`unidade.turnos[0].hora_fim`) para o do formulário plano. */
function nomeNoFormulario(campo: string): string {
  const semIndice = campo.replaceAll(/\[(\d+)\]/g, '.$1');
  if (semIndice.startsWith('unidade.endereco.')) return semIndice.replace('unidade.endereco.', 'unidade_');
  if (semIndice.startsWith('unidade.turnos')) return semIndice.replace('unidade.', '');
  if (semIndice === 'unidade.turnos') return 'turnos';
  return semIndice.replace('responsavel.', 'responsavel_').replace('unidade.', 'unidade_');
}

export function TelaCriarConta() {
  const { t, i18n } = useTranslation();
  const criar = useCriarEmpresa();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [criada, setCriada] = useState<{ email: string; conviteEnviado: boolean } | null>(null);

  const [paisInicial] = useState<Pais>(() => paisDoNavegador(globalThis.navigator?.languages ?? []) ?? 'BR');
  const formulario = useFormulario(esquema, {
    defaultValues: {
      razao_social: '',
      nome_fantasia: '',
      pais: paisInicial,
      identificacao: '',
      // O idioma da interface (escolha guardada ou do navegador) é o ponto de partida; trocar o país o acompanha.
      idioma_padrao: (IDIOMAS as readonly string[]).includes(i18n.language)
        ? (i18n.language as Idioma)
        : PERFIS_DE_PAIS[paisInicial].idioma_padrao,
      fuso_horario: PERFIS_DE_PAIS[paisInicial].fuso_padrao,
      plano: '',
      srv_habilitado: false,
      responsavel_nome: '',
      responsavel_email: '',
      responsavel_telefone: '',
      unidade_codigo: '',
      unidade_nome: '',
      unidade_subdivisao: '',
      unidade_codigo_postal: '',
      turnos: [],
    },
  });
  const { errors, dirtyFields } = formulario.formState;
  const pais = formulario.watch('pais');
  const tipo = PERFIS_DE_PAIS[pais].tipo_empresa;

  /** O país define o tipo de identificação, o fuso e, se a pessoa não escolheu o idioma, o idioma padrão. */
  const aoTrocarOPais = (novo: Pais) => {
    const novoPerfil = PERFIS_DE_PAIS[novo];
    formulario.setValue('fuso_horario', novoPerfil.fuso_padrao);
    if (!dirtyFields.idioma_padrao) formulario.setValue('idioma_padrao', novoPerfil.idioma_padrao);
    formulario.setValue('identificacao', '');
    formulario.clearErrors('identificacao');
    // A subdivisão e o código postal são do país: trocar o país os limpa.
    formulario.setValue('unidade_subdivisao', '');
    formulario.setValue('unidade_codigo_postal', '');
    formulario.clearErrors(['unidade_subdivisao', 'unidade_codigo_postal']);
  };

  // O resumo de erros diz de qual campo é cada mensagem: "Informe o nome" aparece no responsável e na unidade.
  const rotulosDoResumo: Record<string, string> = {
    razao_social: t('empresa.razaoSocial'),
    identificacao: t(`identificacao.tipos.${tipo}`),
    plano: t('empresa.plano'),
    responsavel_nome: t('criarConta.responsavelNome'),
    responsavel_email: t('criarConta.responsavelEmail'),
    responsavel_telefone: t('criarConta.responsavelTelefone'),
    unidade_codigo: `${t('criarConta.secaoUnidade')}: ${t('unidade.codigo')}`,
    unidade_nome: `${t('criarConta.secaoUnidade')}: ${t('unidade.nome')}`,
    unidade_codigo_postal: pais === 'BR' ? t('unidade.cep') : t('unidade.codigoPostal'),
  };

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    const corpo: NovaEmpresa = {
      pais: v.pais,
      identificacao_tipo: PERFIS_DE_PAIS[v.pais].tipo_empresa,
      identificacao: paraOEnvio(v.identificacao),
      idioma_padrao: v.idioma_padrao,
      razao_social: v.razao_social,
      ...(vazioParaIndefinido(v.nome_fantasia) ? { nome_fantasia: v.nome_fantasia } : {}),
      ...(v.fuso_horario ? { fuso_horario: v.fuso_horario } : {}),
      ...(v.plano ? { plano: v.plano as NonNullable<NovaEmpresa['plano']> } : {}),
      ...(v.srv_habilitado ? { srv_habilitado: true } : {}),
      responsavel: {
        nome: v.responsavel_nome,
        email: v.responsavel_email,
        ...(vazioParaIndefinido(v.responsavel_telefone) ? { telefone: v.responsavel_telefone } : {}),
      },
      unidade: {
        codigo: v.unidade_codigo,
        nome: v.unidade_nome,
        ...(v.unidade_subdivisao || vazioParaIndefinido(v.unidade_codigo_postal)
          ? {
              endereco: {
                ...(v.unidade_subdivisao ? { subdivisao: v.unidade_subdivisao } : {}),
                ...(vazioParaIndefinido(v.unidade_codigo_postal) ? { codigo_postal: normalizarCodigoPostal(v.unidade_codigo_postal) } : {}),
              },
            }
          : {}),
        turnos: turnosParaOCorpo(v.turnos),
      },
    };
    criar.mutate(corpo, {
      onSuccess: (conta) => setCriada({ email: v.responsavel_email, conviteEnviado: conta.convite_enviado }),
      onError: (erro) => {
        if (!aplicarErrosDoServidor(erro, formulario.setError, { transformar: nomeNoFormulario })) {
          setMensagem(mensagemDoErro(erro));
        }
      },
    });
  });

  if (criada) {
    return (
      <CartaoPublico titulo={t('criarConta.titulo')}>
        {criada.conviteEnviado ? (
          <Alerta tipo="sucesso" titulo={t('criarConta.criadaTitulo')}>
            {t('criarConta.criadaTexto', { email: criada.email })}
          </Alerta>
        ) : (
          <Alerta tipo="atencao" titulo={t('criarConta.criadaTitulo')}>
            {t('criarConta.conviteNaoEnviado', { email: criada.email })}
          </Alerta>
        )}
        <Link to="/entrar" className="texto-corpo-p underline">
          {t('criarConta.irParaEntrar')}
        </Link>
      </CartaoPublico>
    );
  }

  return (
    <CartaoPublico titulo={t('criarConta.titulo')} largo>
      <p className="texto-corpo text-texto-secundario">{t('criarConta.texto')}</p>
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      <Formulario erros={errors} rotulos={rotulosDoResumo} onSubmit={(e) => void enviar(e)} className="flex flex-col gap-6">
        <fieldset className="flex flex-col gap-4">
          <legend className="mb-2 texto-h4">{t('criarConta.secaoEmpresa')}</legend>
          <CampoTexto rotulo={t('empresa.razaoSocial')} autoComplete="organization" erro={errors.razao_social} {...formulario.register('razao_social')} />
          <CampoTexto rotulo={t('empresa.nomeFantasia')} opcional erro={errors.nome_fantasia} {...formulario.register('nome_fantasia')} />
          <CampoSelecao
            rotulo={t('pais.rotulo')}
            apoio={t('pais.apoio')}
            erro={errors.pais}
            {...formulario.register('pais', { onChange: (e: ChangeEvent<HTMLSelectElement>) => aoTrocarOPais(e.target.value as Pais) })}
          >
            {PAISES.map((codigo) => (
              <option key={codigo} value={codigo}>
                {t(`pais.nomes.${codigo}`)}
              </option>
            ))}
          </CampoSelecao>
          <CampoTexto
            rotulo={t(`identificacao.tipos.${tipo}`)}
            apoio={t('identificacao.exemplo', { exemplo: t(`identificacao.exemplos.${tipo}`) })}
            inputMode={modoDeEntrada(tipo)}
            autoComplete="off"
            dado
            erro={errors.identificacao}
            {...formulario.register('identificacao', { onChange: aplicarMascara(() => PERFIS_DE_PAIS[formulario.getValues('pais')].tipo_empresa) })}
          />
          <CampoSelecao rotulo={t('empresa.idiomaPadrao')} apoio={t('criarConta.idiomaApoio')} erro={errors.idioma_padrao} {...formulario.register('idioma_padrao')}>
            {IDIOMAS.map((idioma) => (
              <option key={idioma} value={idioma} lang={idioma}>
                {t(`idiomas.${idioma}`)}
              </option>
            ))}
          </CampoSelecao>
          <CampoSelecao key={pais} rotulo={t('empresa.fusoHorario')} erro={errors.fuso_horario} {...formulario.register('fuso_horario')}>
            {fusosParaSeletor(pais).map((fuso) => (
              <option key={fuso} value={fuso}>
                {traduzirOuOriginal('fusos', fuso)}
              </option>
            ))}
          </CampoSelecao>
          <CampoSelecao rotulo={t('empresa.plano')} erro={errors.plano} {...formulario.register('plano')}>
            <option value="">{t('criarConta.semPlano')}</option>
            {PLANOS.map((plano) => (
              <option key={plano} value={plano}>
                {t(`planos.${plano}`)}
              </option>
            ))}
          </CampoSelecao>
          <CampoCaixa
            rotulo={t('criarConta.srvRotulo')}
            apoio={t('criarConta.srvApoio')}
            {...formulario.register('srv_habilitado')}
          />
        </fieldset>

        <fieldset className="flex flex-col gap-4">
          <legend className="mb-2 texto-h4">{t('criarConta.secaoResponsavel')}</legend>
          <p className="texto-corpo text-texto-secundario">{t('criarConta.responsavelTexto')}</p>
          <CampoTexto rotulo={t('criarConta.responsavelNome')} autoComplete="name" erro={errors.responsavel_nome} {...formulario.register('responsavel_nome')} />
          <CampoTexto rotulo={t('criarConta.responsavelEmail')} type="email" autoComplete="email" erro={errors.responsavel_email} {...formulario.register('responsavel_email')} />
          <CampoTexto rotulo={t('criarConta.responsavelTelefone')} type="tel" autoComplete="tel" opcional erro={errors.responsavel_telefone} {...formulario.register('responsavel_telefone')} />
        </fieldset>

        <fieldset className="flex flex-col gap-4">
          <legend className="mb-2 texto-h4">{t('criarConta.secaoUnidade')}</legend>
          <CampoTexto rotulo={t('unidade.codigo')} apoio={t('unidade.codigoApoio')} dado erro={errors.unidade_codigo} {...formulario.register('unidade_codigo')} />
          <CampoTexto rotulo={t('unidade.nome')} erro={errors.unidade_nome} {...formulario.register('unidade_nome')} />
          <CampoDeSubdivisao key={pais} pais={pais} erro={errors.unidade_subdivisao} {...formulario.register('unidade_subdivisao')} />
          <CampoTexto rotulo={pais === 'BR' ? t('unidade.cep') : t('unidade.codigoPostal')} apoio={t('identificacao.exemplo', { exemplo: t(`unidade.codigoPostalExemplos.${pais}`) })} opcional dado autoComplete="postal-code" erro={errors.unidade_codigo_postal} {...formulario.register('unidade_codigo_postal')} />
          <EditorDeTurnos control={formulario.control} register={formulario.register} errors={errors} />
        </fieldset>

        <div className="flex items-center justify-end gap-2">
          <Link to="/entrar" className="texto-corpo-p underline">
            {t('criarConta.jaTenhoConta')}
          </Link>
          <Botao type="submit" variante="primario" carregando={criar.isPending}>
            {criar.isPending ? t('criarConta.criando') : t('criarConta.botao')}
          </Botao>
        </div>
      </Formulario>
    </CartaoPublico>
  );
}
