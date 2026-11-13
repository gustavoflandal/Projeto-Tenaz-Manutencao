// Escolha dos perfis do usuário (RF-ADM-004 e RF-ADM-005): pelo menos um.
// Os perfis `prestador` e `locatario` pedem o fornecedor ou o locatário
// vinculado, que só existem com os módulos de terceiros e de locação; por isso
// não são oferecidos aqui.

import { useTranslation } from 'react-i18next';

import { CampoCaixa, textoDoErro, type ErroDeCampo } from '../../../componentes/Campo';
import type { Perfil } from '../../../lib/tipos';

const PERFIS_EXTERNOS = new Set(['prestador', 'locatario']);

export function perfisInternos(perfis: readonly Perfil[]): Perfil[] {
  return perfis.filter((p) => !PERFIS_EXTERNOS.has(p.perfil_origem ?? p.codigo));
}

export interface SeletorDePerfisProps {
  readonly perfis: readonly Perfil[];
  readonly valor: readonly string[];
  readonly aoMudar: (valor: string[]) => void;
  readonly erro?: ErroDeCampo | undefined;
  readonly somenteLeitura?: boolean;
}

export function SeletorDePerfis({ perfis, valor, aoMudar, erro, somenteLeitura }: SeletorDePerfisProps) {
  const { t } = useTranslation();
  const mensagem = textoDoErro(erro);
  const alternar = (codigo: string, marcado: boolean) =>
    aoMudar(marcado ? [...valor, codigo] : valor.filter((c) => c !== codigo));

  return (
    <fieldset aria-describedby={mensagem ? 'erro-perfis' : undefined} className="flex flex-col gap-1">
      <legend className="mb-1 texto-rotulo">{t('usuario.perfis')}</legend>
      <p className="texto-legenda text-texto-secundario">{t('usuario.perfisApoio')}</p>
      <div className="grid gap-x-4 md:grid-cols-2">
        {perfisInternos(perfis).map((perfil) => (
          <CampoCaixa
            key={perfil.codigo}
            name={`perfil-${perfil.codigo}`}
            rotulo={perfil.nome}
            checked={valor.includes(perfil.codigo)}
            disabled={somenteLeitura}
            onChange={(e) => alternar(perfil.codigo, e.target.checked)}
          />
        ))}
      </div>
      {mensagem && (
        <p id="erro-perfis" role="alert" className="texto-legenda text-texto">
          {mensagem}
        </p>
      )}
    </fieldset>
  );
}
