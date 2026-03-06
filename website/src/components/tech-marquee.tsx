"use client";

import { cn } from "@/lib/utils";

// A simple marquee without relying on heavy external libraries
export function TechMarquee() {
    const techs = [
        "Node.js", "Express", "TypeScript", "C#", ".NET", "Golang", "Gin", "PostgreSQL", "MongoDB", "Redis",
        "Node.js", "Express", "TypeScript", "C#", ".NET", "Golang", "Gin", "PostgreSQL", "MongoDB", "Redis"
    ];

    return (
        <div className="relative flex w-full flex-col items-center justify-center overflow-hidden bg-background py-8">
            <div className="relative flex w-full max-w-[1200px] overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_10%,black_90%,transparent)]">
                <div className="flex w-max animate-marquee gap-8 py-2">
                    {techs.map((tech, i) => (
                        <div
                            key={i}
                            className="flex items-center justify-center text-lg font-bold text-muted-foreground/40 transition-colors hover:text-foreground/80 md:text-2xl"
                        >
                            {tech}
                        </div>
                    ))}
                    {/* Duplicate set for seamless looping */}
                    {techs.map((tech, i) => (
                        <div
                            key={`dup-${i}`}
                            className="flex items-center justify-center text-lg font-bold text-muted-foreground/40 transition-colors hover:text-foreground/80 md:text-2xl"
                        >
                            {tech}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
