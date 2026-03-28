#!/usr/bin/env node

/\*\*

- Script para commits automáticos seguindo padrões Conventional Commits em pt-BR
- Agrupa mudanças logicamente e cria commits inteligentes
  \*/

import { execSync, spawnSync } from "child*process";
import * as fs from "fs";
import \_ as path from "path";

interface Mudanca {
arquivo: string;
status: string;
}

interface GrupoCommit {
tipo: string;
descricao: string;
detalhes: string[];
arquivos: string[];
}

// Mapeamento de tipos de commit baseado em padrões de arquivo
function determinarTipoCommit(arquivos: string[]): string {
const tipos: Record<string, number> = {
feat: 0,
fix: 0,
refactor: 0,
style: 0,
docs: 0,
chore: 0,
};

for (const arquivo of arquivos) {
// Novas rotas de API são features
if (arquivo.includes("/api/") && arquivo.includes("route.ts")) {
tipos.feat++;
}
// Novas páginas são features
else if (arquivo.includes("/app/") && arquivo.includes("page.tsx")) {
tipos.feat++;
}
// Novos componentes são features
else if (arquivo.includes("/components/") && !arquivo.includes("ui/")) {
tipos.feat++;
}
// Mudanças em middleware/layout são refactor
else if (
arquivo.includes("middleware.ts") ||
arquivo.includes("layout.tsx")
) {
tipos.refactor++;
}
// Componentes UI existentes são refactor
else if (arquivo.includes("/components/ui/")) {
tipos.refactor++;
}
// Documentação
else if (arquivo.includes(".md") || arquivo.includes("README")) {
tipos.docs++;
}
// Configurações e ferramentas
else if (
arquivo.includes("Makefile") ||
arquivo.includes("package.json") ||
arquivo.includes("scripts/")
) {
tipos.chore++;
}
// Outros arquivos TypeScript são refactor
else if (arquivo.includes(".ts") || arquivo.includes(".tsx")) {
tipos.refactor++;
}
}

// Retorna o tipo mais comum, padrão é refactor
const tipoMaisComum = Object.entries(tipos).reduce((a, b) =>
tipos[a[0]] > tipos[b[1]] ? a : b
)[0];
return tipos[tipoMaisComum] > 0 ? tipoMaisComum : "refactor";
}

// Agrupa arquivos por categoria lógica
function agruparMudancas(mudancas: Mudanca[]): GrupoCommit[] {
const grupos: Record<string, GrupoCommit> = {};
const todasMudancas = mudancas; // Guarda referência para uso posterior

for (const mudanca of mudancas) {
const arquivo = mudanca.arquivo;
let categoria = "outros";

    // Categorização
    if (arquivo.includes("middleware.ts") || arquivo.includes("layout.tsx")) {
      categoria = "config";
    } else if (arquivo.includes("/api/") && arquivo.includes("route.ts")) {
      categoria = "api";
    } else if (arquivo.includes("/components/")) {
      categoria = "componentes";
    } else if (arquivo.includes("/app/") && arquivo.includes("page.tsx")) {
      categoria = "paginas";
    } else if (arquivo.includes("Makefile") || arquivo.includes("scripts/")) {
      categoria = "chore";
    } else if (arquivo.includes(".md")) {
      categoria = "docs";
    }

    if (!grupos[categoria]) {
      grupos[categoria] = {
        tipo: determinarTipoCommit([arquivo]),
        descricao: "",
        detalhes: [],
        arquivos: [],
      };
    }

    grupos[categoria].arquivos.push(arquivo);

}

// Gera descrições e detalhes para cada grupo
const gruposFinais: GrupoCommit[] = [];

for (const [categoria, grupo] of Object.entries(grupos)) {
const arquivos = grupo.arquivos;

    // Descrição baseada na categoria
    let descricao = "";
    const detalhes: string[] = [];

    // Verifica se são arquivos novos (A) ou modificados (M)
    const arquivosNovos = arquivos.filter((a) => {
      const mudanca = todasMudancas.find((m) => m.arquivo === a);
      return mudanca?.status.startsWith("A") || mudanca?.status.includes("??");
    });

    if (categoria === "config") {
      if (arquivosNovos.length > 0) {
        descricao = "adiciona configurações e estrutura da aplicação";
      } else {
        descricao = "atualiza configurações e estrutura da aplicação";
      }
      arquivos.forEach((arquivo) => {
        const nomeArquivo = arquivo.split("/").pop() || arquivo;
        if (arquivo.includes("middleware")) {
          detalhes.push(
            "Atualiza middleware de autenticação e redirecionamentos"
          );
        } else if (arquivo.includes("layout")) {
          detalhes.push("Atualiza layout do dashboard");
        } else {
          detalhes.push(`Atualiza ${nomeArquivo}`);
        }
      });
    } else if (categoria === "api") {
      const rotas = arquivos
        .map((a) => {
          const match = a.match(/\/api\/([^/]+)\/([^/]+)/);
          return match ? `${match[1]}/${match[2]}` : null;
        })
        .filter(Boolean) as string[];
      const rotasUnicas = [...new Set(rotas)];

      if (arquivosNovos.length > 0) {
        descricao = `adiciona rotas de API${
          rotasUnicas.length > 0 ? ` (${rotasUnicas.join(", ")})` : ""
        }`;
      } else {
        descricao = `atualiza rotas de API${
          rotasUnicas.length > 0 ? ` (${rotasUnicas.join(", ")})` : ""
        }`;
      }

      rotasUnicas.forEach((rota) => {
        detalhes.push(
          `${
            arquivosNovos.some((a) => a.includes(rota))
              ? "Adiciona"
              : "Atualiza"
          } endpoint /api/${rota}`
        );
      });
    } else if (categoria === "componentes") {
      const componentes = arquivos
        .map((a) => {
          const match = a.match(/\/components\/([^/]+)/);
          return match ? match[1] : null;
        })
        .filter(Boolean) as string[];
      const componentesUnicos = [...new Set(componentes)];

      if (arquivosNovos.length > 0) {
        descricao = `adiciona componentes${
          componentesUnicos.length > 0
            ? ` (${componentesUnicos.join(", ")})`
            : ""
        }`;
      } else {
        descricao = `atualiza componentes${
          componentesUnicos.length > 0
            ? ` (${componentesUnicos.join(", ")})`
            : ""
        }`;
      }

      componentesUnicos.forEach((comp) => {
        detalhes.push(
          `${
            arquivosNovos.some((a) => a.includes(comp))
              ? "Adiciona"
              : "Atualiza"
          } componente ${comp}`
        );
      });
    } else if (categoria === "paginas") {
      const paginas = arquivos
        .map((a) => {
          const match = a.match(/\/app\/[^/]+\/(.+)\/page\.tsx/);
          return match ? match[1] : null;
        })
        .filter(Boolean) as string[];
      const paginasUnicas = [...new Set(paginas)];

      if (arquivosNovos.length > 0) {
        descricao = `adiciona páginas${
          paginasUnicas.length > 0 ? ` (${paginasUnicas.join(", ")})` : ""
        }`;
      } else {
        descricao = `atualiza páginas${
          paginasUnicas.length > 0 ? ` (${paginasUnicas.join(", ")})` : ""
        }`;
      }

      paginasUnicas.forEach((pagina) => {
        detalhes.push(
          `${
            arquivosNovos.some((a) => a.includes(pagina))
              ? "Adiciona"
              : "Atualiza"
          } página ${pagina.replace(/\//g, " > ")}`
        );
      });
    } else if (categoria === "chore") {
      if (arquivosNovos.length > 0) {
        descricao = "adiciona ferramentas e configurações";
      } else {
        descricao = "atualiza ferramentas e configurações";
      }
      arquivos.forEach((arquivo) => {
        const nomeArquivo = arquivo.split("/").pop() || arquivo;
        if (arquivo.includes("Makefile")) {
          detalhes.push("Adiciona comando para commits automáticos");
        } else if (arquivo.includes("scripts/commits-automaticos")) {
          detalhes.push("Adiciona script de commits automáticos");
        } else {
          detalhes.push(`Atualiza ${nomeArquivo}`);
        }
      });
    } else if (categoria === "docs") {
      descricao =
        arquivosNovos.length > 0
          ? "adiciona documentação"
          : "atualiza documentação";
      arquivos.forEach((arquivo) => {
        detalhes.push(
          `${
            arquivosNovos.includes(arquivo) ? "Adiciona" : "Atualiza"
          } ${arquivo.split("/").pop()}`
        );
      });
    } else {
      descricao =
        arquivosNovos.length > 0
          ? "adiciona arquivos diversos"
          : "atualiza arquivos diversos";
      arquivos.forEach((arquivo) => {
        detalhes.push(
          `${
            arquivosNovos.includes(arquivo) ? "Adiciona" : "Atualiza"
          } ${arquivo.split("/").pop()}`
        );
      });
    }

    gruposFinais.push({
      tipo: grupo.tipo,
      descricao,
      detalhes,
      arquivos,
    });

}

return gruposFinais;
}

// Executa comando git
function executarGit(comando: string): string {
try {
return execSync(comando, { encoding: "utf-8", stdio: "pipe" }).trim();
} catch (erro: unknown) {
const mensagem = erro instanceof Error ? erro.message : String(erro);
throw new Error(`Erro ao executar git: ${mensagem}`);
}
}

// Executa git add de forma segura (evita problemas com caracteres especiais)
function executarGitAdd(arquivo: string): void {
try {
// Usa spawnSync para passar o arquivo como argumento separado, evitando problemas com shell
const resultado = spawnSync("git", ["add", "--", arquivo], {
encoding: "utf-8",
stdio: "pipe",
});

    if (resultado.error) {
      throw resultado.error;
    }

    if (resultado.status !== 0) {
      const erro =
        resultado.stderr?.toString() ||
        resultado.stdout?.toString() ||
        "Erro desconhecido";
      throw new Error(erro);
    }

} catch (erro: unknown) {
const mensagem = erro instanceof Error ? erro.message : String(erro);
throw new Error(`Erro ao adicionar arquivo "${arquivo}": ${mensagem}`);
}
}

// Determina o tipo de branch baseado nas mudanças
function determinarTipoBranch(mudancas: Mudanca[]): string {
const tipos: Record<string, number> = {
feat: 0,
fix: 0,
refactor: 0,
chore: 0,
docs: 0,
style: 0,
};

// Conta arquivos novos vs modificados
const arquivosNovos = mudancas.filter(
(m) => m.status.startsWith("A") || m.status.includes("??")
).length;
const arquivosModificados = mudancas.filter((m) =>
m.status.includes("M")
).length;

for (const mudanca of mudancas) {
const arquivo = mudanca.arquivo;
const ehNovo =
mudanca.status.startsWith("A") || mudanca.status.includes("??");

    // Novas rotas de API são features
    if (arquivo.includes("/api/") && arquivo.includes("route.ts")) {
      tipos.feat += ehNovo ? 3 : 1; // Novos arquivos têm peso maior
    }
    // Novas páginas são features
    else if (arquivo.includes("/app/") && arquivo.includes("page.tsx")) {
      tipos.feat += ehNovo ? 3 : 1;
    }
    // Novos componentes são features
    else if (arquivo.includes("/components/") && !arquivo.includes("ui/")) {
      tipos.feat += ehNovo ? 2 : 1;
    }
    // Correções em rotas/páginas podem ser fixes (se for modificação)
    else if (
      !ehNovo &&
      (arquivo.includes("/api/") || arquivo.includes("/app/")) &&
      mudanca.status.includes("M")
    ) {
      // Se há muitos arquivos modificados, pode ser fix
      if (arquivosModificados > arquivosNovos) {
        tipos.fix += 2;
      } else {
        tipos.feat += 1;
      }
    }
    // Mudanças em middleware/layout são refactor
    else if (
      arquivo.includes("middleware.ts") ||
      arquivo.includes("layout.tsx")
    ) {
      tipos.refactor += 2;
    }
    // Componentes UI existentes são refactor
    else if (arquivo.includes("/components/ui/")) {
      tipos.refactor += 1;
    }
    // Documentação
    else if (arquivo.includes(".md") || arquivo.includes("README")) {
      tipos.docs += ehNovo ? 2 : 1;
    }
    // Configurações e ferramentas
    else if (
      arquivo.includes("Makefile") ||
      arquivo.includes("package.json") ||
      arquivo.includes("scripts/") ||
      arquivo.includes(".gitignore") ||
      arquivo.includes("tsconfig.json") ||
      arquivo.includes("next.config") ||
      arquivo.includes(".env")
    ) {
      tipos.chore += 2;
    }
    // Outros arquivos TypeScript são refactor
    else if (arquivo.includes(".ts") || arquivo.includes(".tsx")) {
      tipos.refactor += 1;
    }
    // Arquivos de estilo
    else if (arquivo.includes(".css") || arquivo.includes(".scss")) {
      tipos.style += 1;
    }

}

// Retorna o tipo mais comum, padrão é feat
const tipoMaisComum = Object.entries(tipos).reduce((a, b) =>
tipos[a[0]] > tipos[b[1]] ? a : b
)[0];
return tipos[tipoMaisComum] > 0 ? tipoMaisComum : "feat";
}

// Busca o próximo número de stage para um tipo de branch
function buscarProximoStage(tipoBranch: string): number {
const stages: number[] = [];
const padrao = new RegExp(`${tipoBranch}/stage-(\\d+)`);

// Busca branches remotas
try {
const branchesRemotas = executarGit("git ls-remote --heads origin");
branchesRemotas.split("\n").forEach((linha) => {
// Formato: <hash> refs/heads/tipo/stage-N
const match = linha.match(
new RegExp(`refs/heads/${tipoBranch}/stage-(\\d+)`)
);
if (match) {
const stageNum = parseInt(match[1], 10);
if (!isNaN(stageNum) && !stages.includes(stageNum)) {
stages.push(stageNum);
}
}
});
} catch (erro: unknown) {
// Se não conseguir buscar branches remotas, continua com branches locais
const mensagem = erro instanceof Error ? erro.message : String(erro);
console.log(`⚠️  Não foi possível buscar branches remotas: ${mensagem}`);
}

// Verifica branches locais
try {
const branchesLocais = executarGit("git branch -a");
branchesLocais.split("\n").forEach((linha) => {
const linhaLimpa = linha
.trim()
.replace(/^\*\s+/, "")
.replace(/^remotes\/[^/]+\//, "");
const match = linhaLimpa.match(padrao);
if (match) {
const stageNum = parseInt(match[1], 10);
if (!isNaN(stageNum) && !stages.includes(stageNum)) {
stages.push(stageNum);
}
}
});
} catch (erro: unknown) {
// Ignora erros ao buscar branches locais
const mensagem = erro instanceof Error ? erro.message : String(erro);
console.log(`⚠️  Não foi possível buscar branches locais: ${mensagem}`);
}

// Retorna o próximo stage (maior + 1, ou 1 se não houver nenhum)
if (stages.length === 0) {
return 1;
}

const maiorStage = Math.max(...stages);
return maiorStage + 1;
}

// Cria branch com tipo e stage sequencial
function criarBranch(mudancas: Mudanca[]): string {
try {
const branchAtual = executarGit("git branch --show-current");

    // Se já está em uma branch que segue o padrão tipo/stage-N, usa ela
    const padraoBranch = /^(feat|fix|refactor|chore|docs|style)\/stage-\d+/;
    if (branchAtual && padraoBranch.test(branchAtual)) {
      console.log(`ℹ️  Usando branch existente: ${branchAtual}`);
      return branchAtual;
    }

    // Determina o tipo de branch baseado nas mudanças
    const tipoBranch = determinarTipoBranch(mudancas);

    // Busca o próximo stage disponível
    const proximoStage = buscarProximoStage(tipoBranch);

    // Cria nome da branch no formato: tipo/stage-N
    const branchNome = `${tipoBranch}/stage-${proximoStage}`;

    // Mostra informações sobre a branch que será criada
    const arquivosNovos = mudancas.filter(
      (m) => m.status.startsWith("A") || m.status.includes("??")
    ).length;
    const arquivosModificados = mudancas.filter((m) =>
      m.status.includes("M")
    ).length;

    console.log(`📋 Tipo de branch: ${tipoBranch}`);
    console.log(`🔢 Stage: ${proximoStage}`);
    console.log(
      `📊 Mudanças: ${arquivosNovos} novo(s), ${arquivosModificados} modificado(s)`
    );

    executarGit(`git checkout -b ${branchNome}`);
    console.log(`✅ Branch criada: ${branchNome}`);
    return branchNome;

} catch (erro: unknown) {
// Se a branch já existe, tenta fazer checkout
const mensagemErro = erro instanceof Error ? erro.message : String(erro);
if (
mensagemErro.includes("already exists") ||
mensagemErro.includes("já existe")
) {
const tipoBranch = determinarTipoBranch(mudancas);
const proximoStage = buscarProximoStage(tipoBranch);
const branchNome = `${tipoBranch}/stage-${proximoStage}`;
executarGit(`git checkout ${branchNome}`);
return branchNome;
}
throw erro;
}
}

// Função principal
function main() {
try {
console.log("🚀 Iniciando commits automáticos...\n");

    // Verifica se há mudanças
    const status = executarGit("git status --porcelain");
    if (!status) {
      console.log("ℹ️  Nenhuma mudança para commitar.");
      return;
    }

    // Obtém lista de arquivos modificados primeiro (para determinar tipo de branch)
    const mudancasRaw = executarGit("git status --porcelain");
    const mudancas: Mudanca[] = mudancasRaw
      .split("\n")
      .filter((linha) => linha.trim())
      .map((linha) => {
        // Formato do git status --porcelain:
        // XY arquivo (onde X = staging, Y = working)
        // ?? arquivo (não rastreado)
        // XY arquivo1 -> arquivo2 (renomeado)

        const status = linha.substring(0, 2).trim();

        // Pega o resto da linha após os 2 primeiros caracteres
        let resto = linha.substring(2).trim();

        // Se for um rename, pega apenas o arquivo de destino
        if (resto.includes(" -> ")) {
          resto = resto.split(" -> ")[1].trim();
        }

        // O arquivo é tudo que sobra após remover o status
        const arquivo = resto;

        return { arquivo, status };
      });

    console.log(`📝 Encontradas ${mudancas.length} mudança(s)\n`);

    // Cria branch baseada nas mudanças
    const branchNome = criarBranch(mudancas);
    console.log("");

    // Executa build antes de fazer commits
    console.log("🔨 Executando build do projeto...");
    try {
      execSync("npm run build", {
        encoding: "utf-8",
        stdio: "inherit",
        cwd: process.cwd(),
      });
      console.log("✅ Build concluído com sucesso!\n");
    } catch (erro: unknown) {
      const mensagem = erro instanceof Error ? erro.message : String(erro);
      console.error("❌ Erro no build:", mensagem);
      console.error(
        "⚠️  O build falhou. Corrija os erros antes de fazer commits."
      );
      process.exit(1);
    }

    // Agrupa mudanças
    const grupos = agruparMudancas(mudancas);

    // Faz commits agrupados
    for (const grupo of grupos) {
      // Adiciona arquivos ao staging de forma segura
      for (const arquivo of grupo.arquivos) {
        executarGitAdd(arquivo);
      }

      // Cria mensagem de commit (escapa caracteres especiais)
      const detalhesFormatados = grupo.detalhes.map((d) => `- ${d}`).join("\n");
      const mensagem = `${grupo.tipo}: ${grupo.descricao}\n\n${detalhesFormatados}`;

      // Faz commit usando arquivo temporário para evitar problemas com caracteres especiais
      const tempFile = path.join(process.cwd(), ".git-commit-msg.txt");
      fs.writeFileSync(tempFile, mensagem, "utf-8");

      try {
        executarGit(`git commit -F "${tempFile}"`);
      } finally {
        // Remove arquivo temporário
        if (fs.existsSync(tempFile)) {
          fs.unlinkSync(tempFile);
        }
      }

      console.log(`✅ Commit criado: ${grupo.tipo}: ${grupo.descricao}`);
      console.log(`   Arquivos: ${grupo.arquivos.length}\n`);
    }

    console.log(`\n✨ Commits concluídos na branch: ${branchNome}`);
    console.log(`📋 Total de commits: ${grupos.length}\n`);

    // Faz push da branch para o remoto
    console.log("📤 Enviando branch para o remoto...");
    try {
      executarGit(`git push -u origin ${branchNome}`);
      console.log(`✅ Branch ${branchNome} enviada para o remoto\n`);
    } catch (erro: unknown) {
      const mensagem = erro instanceof Error ? erro.message : String(erro);
      console.warn(`⚠️  Não foi possível fazer push da branch: ${mensagem}`);
      console.warn("   Você pode fazer o push manualmente depois.\n");
    }

    // Volta para a branch main
    console.log("🔄 Voltando para a branch main...");
    try {
      executarGit("git checkout main");
      console.log("✅ Checkout para main concluído\n");
    } catch (erro: unknown) {
      const mensagem = erro instanceof Error ? erro.message : String(erro);
      console.warn(
        `⚠️  Não foi possível fazer checkout para main: ${mensagem}`
      );
      console.warn("   Você pode fazer o checkout manualmente depois.\n");
    }

    // Executa git pull na main
    console.log("⬇️  Atualizando branch main...");
    try {
      executarGit("git pull origin main");
      console.log("✅ Branch main atualizada\n");
    } catch (erro: unknown) {
      const mensagem = erro instanceof Error ? erro.message : String(erro);
      console.warn(`⚠️  Não foi possível fazer pull na main: ${mensagem}`);
      console.warn("   Você pode fazer o pull manualmente depois.\n");
    }

    console.log("🎉 Processo concluído!");

} catch (erro: unknown) {
const mensagem = erro instanceof Error ? erro.message : String(erro);
console.error("❌ Erro:", mensagem);
process.exit(1);
}
}

main();
