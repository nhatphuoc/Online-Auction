"use client"

import { useState, useEffect } from "react"
import { Navbar } from "@/components/navbar"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Plus, Package, MessageSquare, History } from "lucide-react"
import Link from "next/link"
import { fetchWithAuth } from "@/lib/api-client"

export default function SellingDashboard() {
  const [products, setProducts] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function loadProducts() {
      try {
        const profileRes = await fetchWithAuth("/users/profile/me")
        const profileData = await profileRes.json()
        if (profileData.success) {
          const sellerId = profileData.data.id
          const prodRes = await fetchWithAuth(`/products/seller/${sellerId}`)
          const prodData = await prodRes.json()
          setProducts(prodData)
        }
      } catch (err) {
        console.error("Failed to load seller products", err)
      } finally {
        setLoading(false)
      }
    }
    loadProducts()
  }, [])

  return (
    <div className="min-h-screen bg-black pb-20">
      <Navbar />
      <main className="container mx-auto px-6 pt-32">
        <div className="flex flex-col gap-8">
          <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
            <div>
              <h1 className="text-4xl font-bold tracking-tighter mb-2">Seller Dashboard</h1>
              <p className="text-zinc-500">Manage your active listings and communication.</p>
            </div>
            <Link href="/selling/new">
              <Button className="bg-white text-black hover:bg-zinc-200 rounded-full font-bold px-6 h-12">
                <Plus className="w-5 h-5 mr-2" />
                List New Product
              </Button>
            </Link>
          </div>

          <div className="grid md:grid-cols-3 gap-6">
            {[
              { label: "Active Listings", value: products.length, icon: Package },
              { label: "Unread Messages", value: "3", icon: MessageSquare },
              { label: "Completed Auctions", value: "12", icon: History },
            ].map((stat, i) => (
              <Card key={i} className="bg-white/5 border-white/10 p-6 flex items-center gap-4">
                <div className="w-12 h-12 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-center">
                  <stat.icon className="w-6 h-6 text-zinc-400" />
                </div>
                <div>
                  <p className="text-xs font-bold uppercase tracking-widest text-zinc-500">{stat.label}</p>
                  <p className="text-2xl font-bold text-white">{stat.value}</p>
                </div>
              </Card>
            ))}
          </div>

          <div className="flex flex-col gap-4">
            <h2 className="text-2xl font-bold tracking-tighter">Your Listings</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {loading ? (
                Array(3)
                  .fill(0)
                  .map((_, i) => <div key={i} className="h-64 rounded-3xl bg-white/5 animate-pulse" />)
              ) : products.length > 0 ? (
                products.map((product: any) => (
                  <Card key={product.id} className="bg-white/5 border-white/10 overflow-hidden group">
                    <div className="aspect-video relative overflow-hidden">
                      <img
                        src={product.thumbnailUrl || "/placeholder.svg"}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                        alt={product.name}
                      />
                      <Badge className="absolute top-3 left-3 bg-white text-black font-bold">ACTIVE</Badge>
                    </div>
                    <CardContent className="p-4 flex flex-col gap-4">
                      <div>
                        <h3 className="font-bold text-lg line-clamp-1">{product.name}</h3>
                        <p className="text-zinc-500 text-sm">
                          Current:{" "}
                          <span className="text-white font-bold">${product.currentPrice.toLocaleString()}</span>
                        </p>
                      </div>
                      <div className="flex gap-2">
                        <Link href={`/products/${product.id}`} className="flex-1">
                          <Button
                            variant="outline"
                            className="w-full rounded-full border-white/10 hover:bg-white/5 bg-transparent"
                          >
                            View
                          </Button>
                        </Link>
                        <Button variant="secondary" className="flex-1 rounded-full bg-white/10 text-white">
                          Edit
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ))
              ) : (
                <div className="col-span-full py-20 text-center border-2 border-dashed border-white/5 rounded-3xl">
                  <p className="text-zinc-500">You haven't listed any products yet.</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
