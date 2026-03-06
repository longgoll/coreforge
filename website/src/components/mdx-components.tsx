import type { MDXComponents } from 'mdx/types'
import { ComponentPropsWithoutRef } from 'react'
import { MermaidDiagram } from './mermaid-diagram'
import { Tabs, TabsList, TabsTrigger, TabsContent } from './ui/tabs'

export function useMDXComponents(components: MDXComponents): MDXComponents {
    return {
        h2: (props: ComponentPropsWithoutRef<'h2'>) => {
            const id = typeof props.children === 'string' ? props.children.toLowerCase().replace(/\s+/g, '-') : undefined;
            return <h2 id={id} className="scroll-m-20 border-b border-white/10 pb-2 text-3xl font-semibold tracking-tight transition-colors mt-10 mb-4" {...props} />;
        },
        h3: (props: ComponentPropsWithoutRef<'h3'>) => {
            const id = typeof props.children === 'string' ? props.children.toLowerCase().replace(/\s+/g, '-') : undefined;
            return <h3 id={id} className="scroll-m-20 text-2xl font-semibold tracking-tight mt-8 mb-4 text-orange-400" {...props} />;
        },
        h4: (props: ComponentPropsWithoutRef<'h4'>) => (
            <h4 className="scroll-m-20 text-xl font-semibold tracking-tight mt-6 mb-4" {...props} />
        ),
        p: (props: ComponentPropsWithoutRef<'p'>) => (
            <p className="leading-7 [&:not(:first-child)]:mt-6 text-muted-foreground" {...props} />
        ),
        ul: (props: ComponentPropsWithoutRef<'ul'>) => (
            <ul className="my-6 ml-6 list-disc [&>li]:mt-2 text-muted-foreground" {...props} />
        ),
        ol: (props: ComponentPropsWithoutRef<'ol'>) => (
            <ol className="my-6 ml-6 list-decimal [&>li]:mt-2 text-muted-foreground" {...props} />
        ),
        li: (props: ComponentPropsWithoutRef<'li'>) => (
            <li className="leading-7" {...props} />
        ),
        pre: (props: ComponentPropsWithoutRef<'pre'>) => (
            <pre className="mb-4 mt-6 overflow-x-auto rounded-xl border border-white/10 bg-[#0d0d0d] p-5 custom-scrollbar text-[13px] shadow-2xl" {...props} />
        ),
        code: (props: ComponentPropsWithoutRef<'code'>) => {
            // Intercept mermaid blocks
            if (props.className === 'language-mermaid') {
                return <MermaidDiagram chart={props.children as string} />
            }

            // If it has a newline, it's likely a block code (already wrapped in pre) so don't style as inline
            if (typeof props.children === 'string' && props.children.includes('\n')) {
                return <code {...props} className="font-mono text-zinc-300" />
            }
            return <code className="relative rounded bg-muted/40 px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold text-orange-300" {...props} />
        },
        Tabs,
        TabsList,
        TabsTrigger,
        TabsContent,
        ...components,
    }
}
