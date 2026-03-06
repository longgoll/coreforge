import fs from 'fs';
import path from 'path';

// Define paths
const REGISTRY_PATH = path.join(process.cwd(), '..', 'mock-registry');
const MANIFEST_FILE = path.join(REGISTRY_PATH, 'manifest.json');

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
        }
    }
}

export interface RegistryManifest {
    blueprints: Record<string, unknown>;
    components: { [key: string]: RegistryComponent };
}

let manifestCache: RegistryManifest | null = null;

export async function getRegistry(): Promise<RegistryManifest> {
    if (manifestCache) return manifestCache;
    
    try {
        const fileContent = await fs.promises.readFile(MANIFEST_FILE, 'utf-8');
        manifestCache = JSON.parse(fileContent);
        return manifestCache!;
    } catch (error) {
        console.error('Failed to read registry source:', error);
        return { blueprints: {}, components: {} };
    }
}

export async function getComponent(slug: string): Promise<RegistryComponent | null> {
    const registry = await getRegistry();
    return registry.components[slug] || null;
}

export async function getComponentFileContent(fileUrl: string): Promise<string> {
    try {
        // fileUrl is like "./mock-registry/components/nodejs_express/error-handler/errorHandler.js"
        // remove the "./mock-registry/" prefix because we are already relative to coreforge root
        const normalizedPath = fileUrl.replace('./mock-registry/', '');
        const fullPath = path.join(REGISTRY_PATH, normalizedPath);
        return await fs.promises.readFile(fullPath, 'utf8');
    } catch {
        console.error(`Failed to read source file: ${fileUrl}`);
        return `// Error loading file: ${fileUrl}\n// Please make sure the file exists in the registry.`;
    }
}

export async function getComponentDocs(slug: string, locale: string = 'en'): Promise<string | null> {
    const docsDir = path.join(REGISTRY_PATH, 'components', 'docs');
    
    // Directory-based paths
    const dirLocalizedPath = path.join(docsDir, locale, `${slug}.md`);
    const dirFallbackEnglishPath = path.join(docsDir, 'en', `${slug}.md`);
    
    // Legacy file extension paths
    const localizedPath = path.join(docsDir, `${slug}.${locale}.md`);
    const fallbackEnglishPath = path.join(docsDir, `${slug}.en.md`);
    const legacyPath = path.join(docsDir, `${slug}.md`);

    try {
        if (fs.existsSync(dirLocalizedPath)) {
            return await fs.promises.readFile(dirLocalizedPath, 'utf8');
        }
        if (fs.existsSync(dirFallbackEnglishPath)) {
            return await fs.promises.readFile(dirFallbackEnglishPath, 'utf8');
        }
        if (fs.existsSync(localizedPath)) {
            return await fs.promises.readFile(localizedPath, 'utf8');
        }
        if (fs.existsSync(fallbackEnglishPath)) {
            return await fs.promises.readFile(fallbackEnglishPath, 'utf8');
        }
        if (fs.existsSync(legacyPath)) {
            return await fs.promises.readFile(legacyPath, 'utf8');
        }
    } catch (error) {
        console.error(`Error reading docs for ${slug}:`, error);
    }
    return null;
}

export async function getGeneralDocs(slug: string, locale: string = 'en'): Promise<string | null> {
    const docsDir = path.join(REGISTRY_PATH, 'docs');
    
    const dirLocalizedPath = path.join(docsDir, locale, `${slug}.md`);
    const dirFallbackEnglishPath = path.join(docsDir, 'en', `${slug}.md`);
    
    const localizedPath = path.join(docsDir, `${slug}.${locale}.md`);
    const fallbackEnglishPath = path.join(docsDir, `${slug}.en.md`);

    try {
        if (fs.existsSync(dirLocalizedPath)) {
            return await fs.promises.readFile(dirLocalizedPath, 'utf8');
        }
        if (fs.existsSync(dirFallbackEnglishPath)) {
            return await fs.promises.readFile(dirFallbackEnglishPath, 'utf8');
        }
        if (fs.existsSync(localizedPath)) {
            return await fs.promises.readFile(localizedPath, 'utf8');
        }
        if (fs.existsSync(fallbackEnglishPath)) {
            return await fs.promises.readFile(fallbackEnglishPath, 'utf8');
        }
    } catch (error) {
        console.error(`Error reading general docs for ${slug}:`, error);
    }
    return null;
}
