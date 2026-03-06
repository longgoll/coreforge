"use client"

import React, { useEffect, useRef, useState } from 'react'
import mermaid from 'mermaid'

// Simple hash function for generating safe ids
function hashString(str: string) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash; // Convert to 32bit integer
    }
    return Math.abs(hash).toString(36);
}

export function MermaidDiagram({ chart }: { chart: string }) {
    const [svgContent, setSvgContent] = useState<string>('')
    const elementRef = useRef<HTMLDivElement>(null)

    // Generate a unique ID for this diagram based on its content
    const id = `mermaid-${hashString(chart)}`

    useEffect(() => {
        mermaid.initialize({
            startOnLoad: false,
            theme: 'dark',
            securityLevel: 'loose',
            fontFamily: 'var(--font-geist-sans), sans-serif',
        })

        // Use mermaid.render which returns an object containing the svg string
        mermaid.render(id, chart).then((result) => {
            setSvgContent(result.svg);
        }).catch((err) => {
            console.error('Mermaid render error', err);
            setSvgContent(`<div class="text-red-500">Error rendering diagram</div>`);
        });
    }, [chart, id])

    return (
        <div
            ref={elementRef}
            className="mermaid overflow-x-auto rounded-xl border border-white/5 bg-[#0d0d0d] p-6 shadow-2xl flex justify-center mb-6 mt-4 opacity-90 hover:opacity-100 transition-opacity"
            dangerouslySetInnerHTML={{ __html: svgContent }}
        />
    )
}
