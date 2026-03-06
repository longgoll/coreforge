"use client"

import { useEffect, useState } from "react"

export function TableOfContents() {
    const [headings, setHeadings] = useState<{ id: string; text: string; top: number }[]>([])
    const [activeId, setActiveId] = useState<string>("")

    useEffect(() => {
        const elements = Array.from(document.querySelectorAll("h2, h3"))
            .filter((element) => element.id)
            .map((element) => ({
                id: element.id,
                text: element.textContent ?? "",
                top: (element as HTMLElement).offsetTop,
            }))
        setHeadings(elements)

        const handleScroll = () => {
            const scrollPosition = window.scrollY + 100
            let currentId = ""
            for (const item of elements) {
                if (scrollPosition >= item.top) {
                    currentId = item.id
                } else {
                    break
                }
            }
            if (currentId !== activeId) {
                setActiveId(currentId)
            }
        }

        window.addEventListener("scroll", handleScroll)
        return () => window.removeEventListener("scroll", handleScroll)
    }, [activeId])

    if (headings.length === 0) return null

    return (
        <div className="space-y-4">
            <h4 className="text-sm font-semibold mb-4">On This Page</h4>
            <div className="flex flex-col space-y-2.5 text-sm text-muted-foreground">
                {headings.map((heading) => (
                    <a
                        key={heading.id}
                        href={`#${heading.id}`}
                        className={`hover:text-foreground transition-colors ${activeId === heading.id ? "text-orange-400 font-medium" : ""}`}
                        onClick={(e) => {
                            e.preventDefault();
                            document.querySelector(`#${heading.id}`)?.scrollIntoView({
                                behavior: "smooth"
                            });
                        }}
                    >
                        {heading.text}
                    </a>
                ))}
            </div>
        </div>
    )
}
