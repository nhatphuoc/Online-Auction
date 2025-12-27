"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

export function Navbar() {
  const pathname = usePathname()

  return (
    <nav className="fixed top-0 w-full z-50 border-b border-white/10 bg-black/50 backdrop-blur-md">
      <div className="container mx-auto px-6 h-16 flex items-center justify-between">
        <div className="flex items-center gap-8">
          <Link href="/" className="text-xl font-bold tracking-tighter flex items-center gap-2">
            <div className="w-6 h-6 bg-white rounded-full" />
            ALPACA
          </Link>

          <div className="hidden md:flex items-center gap-6">
            {["Products", "Categories", "Watchlist", "Selling"].map((item) => (
              <Link
                key={item}
                href={`/${item.toLowerCase()}`}
                className={cn(
                  "text-sm font-medium transition-colors hover:text-white",
                  pathname.startsWith(`/${item.toLowerCase()}`) ? "text-white" : "text-zinc-400",
                )}
              >
                {item}
              </Link>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-4">
          <Link href="/sign-in">
            <Button variant="ghost" className="text-zinc-400 hover:text-white hover:bg-white/5">
              Login
            </Button>
          </Link>
          <Link href="/register">
            <Button className="bg-white text-black hover:bg-zinc-200 rounded-full px-6">Sign Up</Button>
          </Link>
        </div>
      </div>
    </nav>
  )
}
