"use client";

import * as React from "react";
import { Copy, Check, Terminal, FolderOpen, FileCode2, Package, Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";

interface CodeBlock {
    name: string;
    language: string;
    content: string;
    highlightedHtml?: string;
}

interface ComponentViewerProps {
    slug: string;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    componentData: any;
    implementations: {
        nodejs: CodeBlock[];
        csharp: CodeBlock[];
        golang: CodeBlock[];
    }
}

function CodePre({ content }: { content: string }) {
    return (
        <pre className="font-mono text-zinc-300 whitespace-pre-wrap"><code>{content}</code></pre>
    );
}

export function ComponentViewer({ slug, componentData, implementations }: ComponentViewerProps) {
    const [copiedContent, setCopiedContent] = React.useState<string | null>(null);
    const [copiedCmd, setCopiedCmd] = React.useState(false);

    const handleCopy = (content: string) => {
        navigator.clipboard.writeText(content);
        setCopiedContent(content);
        setTimeout(() => setCopiedContent(null), 2000);
    };

    const handleCopyCmd = () => {
        navigator.clipboard.writeText(`forge add ${slug}`);
        setCopiedCmd(true);
        setTimeout(() => setCopiedCmd(false), 2000);
    };

    const implKeys = {
        nodejs: "nodejs_express",
        csharp: "csharp_dotnet-webapi",
        golang: "golang_gin",
    } as const;

    return (
        <div className="flex flex-col gap-10 mt-6">

            {/* 1. Installation Block */}
            <section className="flex flex-col gap-3">
                <h2 className="text-2xl font-bold tracking-tight">Installation</h2>
                <div className="relative rounded-xl border bg-black/80 backdrop-blur-md px-5 py-4 font-mono text-sm max-w-full flex items-center justify-between overflow-hidden shadow-lg border-white/10">
                    <div className="flex items-center gap-3">
                        <Terminal className="h-5 w-5 text-orange-400" />
                        <span className="text-zinc-200">forge add {slug}</span>
                    </div>
                    <Button variant="ghost" size="icon" className="h-9 w-9 text-zinc-400 hover:text-white hover:bg-white/10 transition-colors" onClick={handleCopyCmd}>
                        {copiedCmd ? <Check className="h-4 w-4 text-emerald-400" /> : <Copy className="h-4 w-4" />}
                    </Button>
                </div>
            </section>

            {/* Main Tabs for Language Specific Details */}
            <Tabs defaultValue="nodejs" className="relative w-full">
                <div className="flex items-center justify-between pb-4 border-b border-border/50 mb-6">
                    <TabsList className="bg-muted/50 p-1">
                        <TabsTrigger value="nodejs" className="font-medium">Node.js</TabsTrigger>
                        <TabsTrigger value="csharp" className="font-medium">C# .NET</TabsTrigger>
                        <TabsTrigger value="golang" className="font-medium">Go (Gin)</TabsTrigger>
                    </TabsList>
                </div>

                {['nodejs', 'csharp', 'golang'].map(langKey => {
                    const files = implementations[langKey as keyof typeof implementations];
                    const meta = componentData.implementations[implKeys[langKey as keyof typeof implKeys]];

                    return (
                        <TabsContent key={langKey} value={langKey} className="mt-0 flex flex-col gap-10 animate-in fade-in slide-in-from-bottom-2 duration-300">

                            {/* 2. What you get (File Tree) */}
                            {files.length > 0 && (
                                <section className="flex flex-col gap-4">
                                    <div className="flex items-center gap-2">
                                        <FolderOpen className="w-5 h-5 text-muted-foreground" />
                                        <h3 className="text-xl font-bold tracking-tight">What you get</h3>
                                    </div>
                                    <div className="rounded-xl border bg-card p-5">
                                        <div className="flex flex-col gap-2 font-mono text-sm">
                                            {files.map(f => (
                                                <div key={f.name} className="flex items-center gap-2 text-muted-foreground">
                                                    <FileCode2 className="w-4 h-4 text-orange-400/80" />
                                                    <span className="text-foreground/90">{f.name}</span>
                                                    <Badge variant="secondary" className="px-1.5 py-0 text-[10px] h-4 bg-orange-500/10 text-orange-400 border-none ml-2">New</Badge>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </section>
                            )}

                            {/* 3. Dependencies X-Ray */}
                            {meta?.dependencies && meta.dependencies.length > 0 && (
                                <section className="flex flex-col gap-4">
                                    <div className="flex items-center gap-2">
                                        <Package className="w-5 h-5 text-muted-foreground" />
                                        <h3 className="text-xl font-bold tracking-tight">Dependencies</h3>
                                    </div>
                                    <div className="rounded-xl border border-blue-500/20 bg-blue-500/5 p-5">
                                        <div className="flex items-start gap-3">
                                            <Info className="w-5 h-5 text-blue-400 mt-0.5 shrink-0" />
                                            <div>
                                                <p className="text-sm text-muted-foreground mb-3 leading-relaxed">
                                                    The CLI will automatically install these packages into your project if they are missing.
                                                </p>
                                                <div className="flex flex-wrap gap-2">
                                                    {meta.dependencies.map((dep: string) => (
                                                        <code key={dep} className="px-2 py-1 rounded bg-black/40 text-blue-300 font-mono text-sm border border-blue-500/20">
                                                            {dep}
                                                        </code>
                                                    ))}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </section>
                            )}

                            {/* 4. Code Blocks */}
                            <section className="flex flex-col gap-4">
                                <h3 className="text-xl font-bold tracking-tight">Source Code</h3>
                                {files.length === 0 ? (
                                    <div className="rounded-xl border bg-muted/20 p-8 text-center text-muted-foreground text-sm">
                                        Pending implementation for {langKey}.
                                    </div>
                                ) : (
                                    <div className="flex flex-col gap-6">
                                        {files.map(file => (
                                            <div key={file.name} className="relative group rounded-xl overflow-hidden border border-white/10 bg-[#0d0d0d] shadow-2xl">
                                                <div className="flex items-center justify-between px-4 py-3 bg-[#161616] border-b border-white/5 sticky top-0 z-20">
                                                    <div className="flex items-center gap-3">
                                                        <div className="flex space-x-1.5">
                                                            <div className="w-3 h-3 rounded-full bg-red-500/80 border border-black/20"></div>
                                                            <div className="w-3 h-3 rounded-full bg-yellow-500/80 border border-black/20"></div>
                                                            <div className="w-3 h-3 rounded-full bg-green-500/80 border border-black/20"></div>
                                                        </div>
                                                        <div className="text-xs font-mono text-zinc-400 font-medium">{file.name}</div>
                                                    </div>
                                                    <Button variant="ghost" size="icon" className="h-8 w-8 text-zinc-400 hover:text-white transition-colors hover:bg-white/10" onClick={() => handleCopy(file.content)}>
                                                        {copiedContent === file.content ? <Check className="h-4 w-4 text-emerald-400" /> : <Copy className="h-4 w-4" />}
                                                    </Button>
                                                </div>
                                                <div className="overflow-x-auto text-[13px] leading-relaxed max-h-[600px] overflow-y-auto w-full custom-scrollbar relative z-10">
                                                    {file.highlightedHtml ? (
                                                        <div
                                                            className="[&>pre]:!bg-transparent [&>pre]:p-5 [&>pre]:m-0"
                                                            dangerouslySetInnerHTML={{ __html: file.highlightedHtml }}
                                                        />
                                                    ) : (
                                                        <CodePre content={file.content} />
                                                    )}
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </section>

                            {/* 5. Post Install / Usage */}
                            {meta?.postInstall && (
                                <section className="flex flex-col gap-4 mt-4">
                                    <h3 className="text-xl font-bold tracking-tight">Post Installation</h3>
                                    <div className="rounded-xl border bg-card p-5">
                                        <div className="text-sm text-muted-foreground whitespace-pre-wrap leading-relaxed">
                                            {meta.postInstall}
                                        </div>
                                    </div>
                                </section>
                            )}

                        </TabsContent>
                    )
                })}
            </Tabs>
        </div>
    );
}
