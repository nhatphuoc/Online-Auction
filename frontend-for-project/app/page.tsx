import { Navbar } from "@/components/navbar"
import { Button } from "@/components/ui/button"

export default function HomePage() {
  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-6">
        <div className="container mx-auto">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            <div>
              <h1 className="text-6xl md:text-8xl font-bold tracking-tighter leading-[0.9] mb-8">
                The complete platform to <span className="text-zinc-500">bid the web.</span>
              </h1>
              <p className="text-zinc-400 text-lg md:text-xl max-w-md mb-10 leading-relaxed">
                Your team's toolkit to stop configuring and start winning. Securely bid, sell, and scale your auction
                experience with Alpaca.
              </p>
              <div className="flex flex-wrap gap-4">
                <Button
                  size="lg"
                  className="bg-white text-black hover:bg-zinc-200 rounded-full px-8 h-14 text-base font-semibold"
                >
                  Start Bidding
                </Button>
                <Button
                  size="lg"
                  variant="outline"
                  className="border-zinc-800 text-white hover:bg-white/5 rounded-full px-8 h-14 text-base font-semibold bg-transparent"
                >
                  Explore Products
                </Button>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="border-t border-white/10">
        <div className="container mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-4 divide-y md:divide-y-0 md:divide-x divide-white/10">
            {[
              { label: "saved on daily builds.", value: "20 days", sub: "NETFLIX" },
              { label: "faster time to market.", value: "98% faster", sub: "Tripadvisor" },
              { label: "increase in SEO.", value: "300% increase", sub: "box" },
              { label: "faster to build + deploy.", value: "6x faster", sub: "ebay" },
            ].map((stat, i) => (
              <div key={i} className="p-8 md:p-12 group hover:bg-white/[0.02] transition-colors">
                <p className="text-zinc-400 text-sm mb-1 group-hover:text-zinc-300 transition-colors">
                  <span className="text-white font-bold">{stat.value}</span> {stat.label}
                </p>
                <h3 className="text-2xl font-bold tracking-tighter opacity-50 group-hover:opacity-100 transition-opacity">
                  {stat.sub.toUpperCase()}
                </h3>
              </div>
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}
