"use client"

import { useState, useEffect } from "react"
import { Navbar } from "@/components/navbar"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Users, Gavel, FileText, CheckCircle, XCircle } from "lucide-react"
import { fetchWithAuth } from "@/lib/api-client"

export default function AdminDashboard() {
  const [requests, setRequests] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function loadRequests() {
      try {
        const res = await fetchWithAuth("/users")
        const data = await res.json()
        setRequests(data.content || [])
      } catch (err) {
        console.error("Failed to load upgrade requests", err)
      } finally {
        setLoading(false)
      }
    }
    loadRequests()
  }, [])

  const approveRequest = async (id: number) => {
    try {
      const res = await fetchWithAuth(`/users/${id}/approve`, { method: "POST" })
      if (res.ok) {
        setRequests((prev) => prev.filter((r: any) => r.id !== id))
        alert("User upgraded to SELLER")
      }
    } catch (err) {
      alert("Approval failed")
    }
  }

  return (
    <div className="min-h-screen bg-black pb-20">
      <Navbar />
      <main className="container mx-auto px-6 pt-32">
        <div className="flex flex-col gap-8">
          <div>
            <h1 className="text-4xl font-bold tracking-tighter mb-2">Admin Command Center</h1>
            <p className="text-zinc-500">Platform-wide management and analytics.</p>
          </div>

          <div className="grid md:grid-cols-4 gap-6">
            {[
              { label: "Total Users", value: "1,240", icon: Users },
              { label: "Active Auctions", value: "86", icon: Gavel },
              { label: "Pending Upgrades", value: requests.length, icon: FileText },
              { label: "Revenue (MTD)", value: "$12,400", icon: CheckCircle },
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
            <h2 className="text-2xl font-bold tracking-tighter">Seller Upgrade Requests</h2>
            <div className="flex flex-col gap-3">
              {loading ? (
                Array(3)
                  .fill(0)
                  .map((_, i) => <div key={i} className="h-20 rounded-2xl bg-white/5 animate-pulse" />)
              ) : requests.length > 0 ? (
                requests.map((req: any) => (
                  <div
                    key={req.id}
                    className="flex items-center justify-between p-6 rounded-2xl bg-white/5 border border-white/10"
                  >
                    <div className="flex flex-col gap-1">
                      <div className="flex items-center gap-2">
                        <span className="font-bold text-lg">User #{req.userId}</span>
                        <Badge variant="outline" className="border-white/10 text-zinc-500">
                          PENDING
                        </Badge>
                      </div>
                      <p className="text-sm text-zinc-400 italic">"{req.reason}"</p>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-zinc-500 hover:text-white hover:bg-white/5 rounded-full px-4"
                      >
                        <XCircle className="w-4 h-4 mr-2" />
                        Decline
                      </Button>
                      <Button
                        size="sm"
                        onClick={() => approveRequest(req.id)}
                        className="bg-white text-black hover:bg-zinc-200 rounded-full px-6 font-bold"
                      >
                        <CheckCircle className="w-4 h-4 mr-2" />
                        Approve
                      </Button>
                    </div>
                  </div>
                ))
              ) : (
                <div className="py-12 text-center border-2 border-dashed border-white/5 rounded-3xl">
                  <p className="text-zinc-500">No pending upgrade requests.</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
