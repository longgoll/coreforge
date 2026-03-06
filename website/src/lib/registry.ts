export interface RegistryComponent {
  description: string;
  category: string;
  tags: string[];
  implementations: {
    [key: string]: {
      files: { url: string; target: string }[];
      dependencies?: string[];
      envVars?: { key: string; default: string; description: string }[];
      installCmd?: string;
      postInstall?: string;
      requires?: string[];
      conflicts?: string[];
    };
  };
}

export interface RegistryManifest {
  blueprints: Record<string, unknown>;
  components: { [key: string]: RegistryComponent };
}

const REPO_OWNER = 'longgoll';
const REPO_NAME = 'forge-registry';
const BRANCH = 'main';

const BASE_URL = `https://raw.githubusercontent.com/${REPO_OWNER}/${REPO_NAME}/${BRANCH}`;
const MANIFEST_URL = `${BASE_URL}/manifest.json`;

export async function getRegistry(): Promise<RegistryManifest> {
  try {
    const res = await fetch(MANIFEST_URL, { next: { revalidate: 3600 } });
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
    return await res.json();
  } catch (error) {
    console.error("Failed to fetch registry source:", error);
    return { blueprints: {}, components: {} };
  }
}

export async function getComponent(
  slug: string,
): Promise<RegistryComponent | null> {
  const registry = await getRegistry();
  return registry.components[slug] || null;
}

export async function getComponentFileContent(
  fileUrl: string,
): Promise<string> {
  try {
    // If the fileUrl is already a full URL (which manifest should provide)
    // or if it's a relative path starting with './mock-registry/' or './forge-registry/'
    let url = fileUrl;
    if (!fileUrl.startsWith('http')) {
      const normalizedPath = fileUrl.replace(/^\.\/(mock-registry|forge-registry)\//, '');
      url = `${BASE_URL}/${normalizedPath}`;
    }

    const res = await fetch(url, { next: { revalidate: 3600 } });
    if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
    return await res.text();
  } catch (error) {
    console.error(`Failed to fetch source file: ${fileUrl}`, error);
    return `// Error loading file: ${fileUrl}\n// Please make sure the file exists in the registry.`;
  }
}

export async function getComponentDocs(
  slug: string,
  locale: string = "en",
): Promise<string | null> {
  const docsPaths = [
    `components/docs/${locale}/${slug}.md`,
    `components/docs/en/${slug}.md`,
    `components/docs/${slug}.${locale}.md`,
    `components/docs/${slug}.en.md`,
    `components/docs/${slug}.md`
  ];

  for (const docPath of docsPaths) {
    try {
      const url = `${BASE_URL}/${docPath}`;
      const res = await fetch(url, { next: { revalidate: 3600 } });
      if (res.ok) {
        return await res.text();
      }
    } catch (error) {
      // Ignore and continue trying next fallback path
    }
  }
  
  return null;
}

export async function getGeneralDocs(
  slug: string,
  locale: string = "en",
): Promise<string | null> {
  const docsPaths = [
    `docs/${locale}/${slug}.md`,
    `docs/en/${slug}.md`,
    `docs/${slug}.${locale}.md`,
    `docs/${slug}.en.md`
  ];

  for (const docPath of docsPaths) {
    try {
      const url = `${BASE_URL}/${docPath}`;
      const res = await fetch(url, { next: { revalidate: 3600 } });
      if (res.ok) {
        return await res.text();
      }
    } catch (error) {
      // Ignore and continue trying next fallback path
    }
  }
  
  return null;
}
