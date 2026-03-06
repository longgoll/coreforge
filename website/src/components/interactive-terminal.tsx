"use client";

import { useState, useEffect } from "react";
import { Terminal } from "lucide-react";

const commands = [
    { text: "forge add jwt-auth", type: "input", delay: 800 },
    { text: "✓ Found component: jwt-auth (Node.js Express)", type: "success", delay: 1500 },
    { text: "✓ Installed dependencies: jsonwebtoken", type: "success", delay: 2000 },
    { text: "✓ Created src/middlewares/auth.middleware.ts", type: "success", delay: 2200 },
    { text: "✓ Created src/services/token.service.ts", type: "success", delay: 2400 },
    { text: "✨ Done in 1.4s", type: "info", delay: 2800 },
];

export function InteractiveTerminal() {
    const [lines, setLines] = useState<number>(0);
    const [typedChars, setTypedChars] = useState<number>(0);

    useEffect(() => {
        let timeout: NodeJS.Timeout;

        const runSequence = async () => {
            // Typing effect for first command
            for (let i = 0; i <= commands[0].text.length; i++) {
                await new Promise((r) => setTimeout(r, 60)); // typing speed
                setTypedChars(i);
            }

            setLines(1); // after typing, consider line 1 done

            // Execute rest of the lines based on delays
            commands.forEach((cmd, idx) => {
                if (idx === 0) return;
                setTimeout(() => {
                    setLines((prev) => Math.max(prev, idx + 1));
                }, cmd.delay);
            });

            // Reset loop
            timeout = setTimeout(() => {
                setLines(0);
                setTypedChars(0);
            }, commands[commands.length - 1].delay + 4000);
        };

        if (lines === 0) {
            runSequence();
        }

        return () => clearTimeout(timeout);
    }, [lines]);

    return (
        <div className="w-full max-w-lg rounded-xl overflow-hidden border border-white/10 bg-black/80 backdrop-blur-md shadow-2xl">
            {/* Terminal Header */}
            <div className="flex z-10 items-center px-4 py-3 border-b border-white/5 bg-white/5">
                <div className="flex gap-2 mr-4">
                    <div className="w-3 h-3 rounded-full bg-red-500/80"></div>
                    <div className="w-3 h-3 rounded-full bg-yellow-500/80"></div>
                    <div className="w-3 h-3 rounded-full bg-green-500/80"></div>
                </div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground font-medium">
                    <Terminal size={12} />
                    <span>bash — forge-cli</span>
                </div>
            </div>

            {/* Terminal Body */}
            <div className="p-5 font-mono text-sm leading-relaxed h-[240px] overflow-hidden text-left bg-transparent">
                <div className="flex flex-col gap-1.5">
                    {/* Input Line */}
                    <div className="flex items-center text-zinc-300">
                        <span className="text-emerald-400 mr-2">~</span>
                        <span className="text-zinc-500 mr-2">$</span>
                        <span>
                            <span className="text-pink-400">forge</span>{" "}
                            {commands[0].text.substring(6, typedChars)}
                        </span>
                        {lines === 0 && <span className="w-2 h-4 bg-zinc-400 animate-pulse ml-1" />}
                    </div>

                    {/* Output Lines */}
                    {commands.slice(1, lines).map((cmd, idx) => (
                        <div
                            key={idx}
                            className={`animate-in fade-in slide-in-from-bottom-1 duration-300 ${cmd.type === "success" ? "text-emerald-400" : "text-sky-400"
                                }`}
                        >
                            {cmd.text}
                        </div>
                    ))}

                    {/* Prompt ready again */}
                    {lines === commands.length && (
                        <div className="flex items-center text-zinc-300 mt-2">
                            <span className="text-emerald-400 mr-2">~</span>
                            <span className="text-zinc-500 mr-2">$</span>
                            <span className="w-2 h-4 bg-zinc-400 animate-pulse ml-1" />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
