import { ChevronRight } from "lucide-react";
import { Link } from "@/i18n/routing";
import { ComponentViewer } from "@/components/component-viewer";
import { getComponent, getComponentFileContent, getComponentDocs } from "@/lib/registry";
import { notFound } from "next/navigation";
import { MDXRemote } from 'next-mdx-remote/rsc';
import { useMDXComponents } from "@/components/mdx-components";
import { TableOfContents } from "@/components/table-of-contents";
import remarkGfm from 'remark-gfm';
import { codeToHtml } from 'shiki';

export default async function ComponentPage({ params }: { params: Promise<{ locale: string, slug: string }> }) {
    const resolvedParams = await params;
    const { locale, slug } = resolvedParams;

    const componentData = await getComponent(slug);
    if (!componentData) {
        notFound();
    }

    const title = slug.split('-').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');

    // Fetch Markdown Documentation
    const docsMarkdown = await getComponentDocs(slug, locale);

    // Only grab content after the frontmatter (quick hack for basic frontmatter removal)
    let cleanedDocs = docsMarkdown || '';
    if (cleanedDocs.startsWith('---')) {
        const endOfFrontmatter = cleanedDocs.indexOf('---', 3);
        if (endOfFrontmatter !== -1) {
            cleanedDocs = cleanedDocs.substring(endOfFrontmatter + 3).trim();
        }
    }

    const mapFiles = async (implementationKey: string) => {
        const impl = componentData.implementations[implementationKey];
        if (!impl || !impl.files) return [];
        return Promise.all(impl.files.map(async file => {
            const name = file.target.startsWith('/') ? file.target.substring(1) : file.target;
            const language = file.target.endsWith('.cs') ? 'csharp' : file.target.endsWith('.go') ? 'go' : 'typescript';
            const content = await getComponentFileContent(file.url);

            let highlightedHtml = '';
            try {
                highlightedHtml = await codeToHtml(content, {
                    lang: language,
                    theme: 'github-dark',
                });
            } catch (e) {
                console.error("Shiki error:", e);
                highlightedHtml = `<pre class="font-mono text-zinc-300 whitespace-pre-wrap"><code>${content}</code></pre>`;
            }

            return {
                name,
                language,
                content,
                highlightedHtml
            };
        }));
    };

    const implementations = {
        nodejs: await mapFiles('nodejs_express'),
        csharp: await mapFiles('csharp_dotnet-webapi'),
        golang: await mapFiles('golang_gin')
    };

    return (
        <>
            <main className="relative py-6 lg:py-8 w-full max-w-3xl mx-auto md:mx-0">
                <div className="mx-auto w-full min-w-0">
                    <div className="mb-4 flex items-center space-x-1 text-sm text-muted-foreground">
                        <Link className="overflow-hidden text-ellipsis whitespace-nowrap hover:underline" href="/docs">Docs</Link>
                        <ChevronRight className="h-4 w-4" />
                        <Link className="overflow-hidden text-ellipsis whitespace-nowrap hover:underline" href="/components">Components</Link>
                        <ChevronRight className="h-4 w-4" />
                        <div className="font-medium text-foreground">{slug}</div>
                    </div>

                    <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">{title}</h1>
                    <p className="text-xl text-muted-foreground mb-10 max-w-2xl leading-relaxed">
                        {componentData.description}
                    </p>

                    {/* Render the markdown content */}
                    {cleanedDocs && (
                        <div className="mb-14 pb-8 border-b border-white/5">
                            <MDXRemote
                                source={cleanedDocs}
                                components={useMDXComponents({})}
                                options={{ mdxOptions: { remarkPlugins: [remarkGfm] } }}
                            />
                        </div>
                    )}

                    {/* The 3-tab visualizer is injected below the story */}
                    <h2 id="component-source-code" className="scroll-m-20 border-b border-white/10 pb-2 text-3xl font-semibold tracking-tight transition-colors mt-10 mb-4">
                        Component Source Code
                    </h2>
                    <ComponentViewer slug={slug} componentData={componentData} implementations={implementations} />
                </div>
            </main>

            <div className="hidden xl:block">
                <div className="sticky top-20 h-[calc(100vh-5rem)] overflow-y-auto pt-8">
                    <TableOfContents />
                </div>
            </div>
        </>
    );
}
