"use client"

import { useState, useEffect } from "react"
import { useParams } from "next/navigation"
import { Navbar } from "@/components/navbar"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Clock, ShieldCheck, User } from "lucide-react"

export default function ProductDetailPage() {
  const { id } = useParams()
  const [product, setProduct] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [bidAmount, setBidAmount] = useState("")

  useEffect(() => {
    async function fetchProduct() {
      try {
        const res = await fetch(`http://localhost:8080/api/products/${id}`)
        const data = await res.json()
        setProduct(data)
      } catch (error) {
        console.error("Failed to fetch product:", error)
      } finally {
        setLoading(false)
      }
    }
    fetchProduct()
  }, [id])

  if (loading) return <div className="min-h-screen bg-black" />
  if (!product) return <div className="min-h-screen flex items-center justify-center">Product not found</div>

  const minBid = product.currentPrice + product.stepPrice

  return (
    <div className="min-h-screen pb-20">
      <Navbar />

      <main className="container mx-auto px-6 pt-32">
        <div className="grid lg:grid-cols-2 gap-12">
          {/* Image Gallery */}
          <div className="flex flex-col gap-4">
            <div className="aspect-[4/3] rounded-3xl overflow-hidden bg-white/5 border border-white/10">
              <img
                src={product.thumbnailUrl || "/placeholder.svg?height=600&width=800&query=product"}
                alt={product.name}
                className="w-full h-full object-cover"
              />
            </div>
            <div className="grid grid-cols-4 gap-4">
              {product.images?.map((img: string, i: number) => (
                <div
                  key={i}
                  className="aspect-square rounded-xl overflow-hidden bg-white/5 border border-white/10 cursor-pointer hover:border-white/30"
                >
                  <img
                    src={img || "/placeholder.svg"}
                    alt={`${product.name} ${i}`}
                    className="w-full h-full object-cover"
                  />
                </div>
              ))}
            </div>
          </div>

          {/* Product Info & Bidding */}
          <div className="flex flex-col gap-8">
            <div className="flex flex-col gap-4">
              <div className="flex items-center gap-2">
                <Badge variant="secondary" className="bg-white/10 text-zinc-300 rounded-full">
                  {product.categoryName}
                </Badge>
                <div className="flex items-center gap-1.5 text-zinc-400 text-sm ml-2">
                  <Clock className="w-4 h-4" />
                  <span>3 days left</span>
                </div>
              </div>
              <h1 className="text-4xl md:text-5xl font-bold tracking-tighter leading-tight">{product.name}</h1>
            </div>

            <div className="p-6 rounded-3xl bg-white/5 border border-white/10 flex flex-col gap-6">
              <div className="flex justify-between items-end">
                <div className="flex flex-col gap-1">
                  <span className="text-zinc-500 text-sm font-medium">Current Bid</span>
                  <div className="flex items-baseline gap-2">
                    <span className="text-4xl font-bold tracking-tighter">
                      ${product.currentPrice.toLocaleString()}
                    </span>
                    <span className="text-zinc-500 text-sm">approx. €{Math.round(product.currentPrice * 0.92)}</span>
                  </div>
                </div>
                <div className="text-right">
                  <span className="text-zinc-500 text-sm font-medium">Reserve price</span>
                  <div className="flex items-center gap-1 text-white font-semibold">
                    <ShieldCheck className="w-4 h-4 text-emerald-500" />
                    Met
                  </div>
                </div>
              </div>

              <div className="flex flex-col gap-3">
                <div className="flex gap-2">
                  <Input
                    type="number"
                    placeholder={`Min. bid $${minBid.toLocaleString()}`}
                    value={bidAmount}
                    onChange={(e) => setBidAmount(e.target.value)}
                    className="h-14 rounded-full bg-white/5 border-white/10 px-6 text-lg font-bold focus:border-white/30"
                  />
                  <Button className="h-14 px-8 bg-white text-black hover:bg-zinc-200 rounded-full font-bold text-lg flex-1">
                    Place Bid
                  </Button>
                </div>
                <p className="text-zinc-500 text-xs text-center">
                  By placing a bid, you agree to our <span className="underline cursor-pointer">Terms of Service</span>.
                </p>
              </div>

              {product.buyNowPrice && (
                <div className="pt-6 border-t border-white/10">
                  <Button
                    variant="outline"
                    className="w-full h-14 rounded-full border-white/10 hover:bg-white/5 font-bold text-lg bg-transparent"
                  >
                    Buy It Now for ${product.buyNowPrice.toLocaleString()}
                  </Button>
                </div>
              )}
            </div>

            {/* Seller Info */}
            <div className="flex items-center justify-between p-4 rounded-2xl bg-white/5 border border-white/10">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-full bg-zinc-800 flex items-center justify-center">
                  <User className="w-6 h-6 text-zinc-500" />
                </div>
                <div>
                  <p className="text-sm font-bold">{product.sellerInfo?.username || "Trusted Seller"}</p>
                  <p className="text-xs text-zinc-500">98% Positive Feedback (420 reviews)</p>
                </div>
              </div>
              <Button variant="ghost" className="text-zinc-400 hover:text-white">
                View Profile
              </Button>
            </div>
          </div>
        </div>

        {/* Detailed Description & Info Tabs */}
        <div className="mt-20">
          <Tabs defaultValue="description" className="w-full">
            <TabsList className="bg-transparent border-b border-white/10 w-full justify-start rounded-none h-auto p-0 gap-8">
              <TabsTrigger
                value="description"
                className="bg-transparent border-none rounded-none data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:border-b-2 data-[state=active]:border-white px-0 py-4 text-lg font-bold text-zinc-500 transition-all"
              >
                Description
              </TabsTrigger>
              <TabsTrigger
                value="history"
                className="bg-transparent border-none rounded-none data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:border-b-2 data-[state=active]:border-white px-0 py-4 text-lg font-bold text-zinc-500 transition-all"
              >
                Bid History
              </TabsTrigger>
              <TabsTrigger
                value="qa"
                className="bg-transparent border-none rounded-none data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:border-b-2 data-[state=active]:border-white px-0 py-4 text-lg font-bold text-zinc-500 transition-all"
              >
                Questions
              </TabsTrigger>
            </TabsList>
            <TabsContent value="description" className="pt-10">
              <div className="prose prose-invert max-w-none">
                <p className="text-zinc-400 text-lg leading-relaxed whitespace-pre-wrap">{product.description}</p>
              </div>
            </TabsContent>
            <TabsContent value="history" className="pt-10">
              <div className="flex flex-col gap-4">
                {/* Mock bid history for now */}
                {[
                  { user: "j***e", amount: 25000000, time: "2 hours ago" },
                  { user: "a***s", amount: 24500000, time: "5 hours ago" },
                  { user: "m***k", amount: 24000000, time: "1 day ago" },
                ].map((bid, i) => (
                  <div
                    key={i}
                    className="flex justify-between items-center p-4 rounded-xl bg-white/5 border border-white/10"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-zinc-800 flex items-center justify-center">
                        <User className="w-4 h-4 text-zinc-500" />
                      </div>
                      <span className="font-bold text-zinc-300">{bid.user}</span>
                    </div>
                    <div className="text-right">
                      <p className="font-bold text-white">${bid.amount.toLocaleString()}</p>
                      <p className="text-xs text-zinc-500">{bid.time}</p>
                    </div>
                  </div>
                ))}
              </div>
            </TabsContent>
            <TabsContent value="qa" className="pt-10">
              <div className="flex flex-col gap-6">
                <div className="p-6 rounded-3xl bg-white/5 border border-white/10 flex flex-col gap-4">
                  <h3 className="text-xl font-bold">Ask the seller</h3>
                  <div className="flex gap-4">
                    <Input
                      placeholder="Type your question here..."
                      className="h-12 rounded-full bg-white/5 border-white/10 px-6 focus:border-white/30"
                    />
                    <Button className="h-12 rounded-full bg-white text-black hover:bg-zinc-200 px-8 font-bold">
                      Send
                    </Button>
                  </div>
                </div>
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </main>
    </div>
  )
}
