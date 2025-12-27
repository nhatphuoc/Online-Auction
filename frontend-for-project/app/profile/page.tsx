"use client"

import { useState, useEffect } from "react"
import { Navbar } from "@/components/navbar"
import { authService } from "@/services/auth-service"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { User, Mail, Phone, Calendar, Shield } from "lucide-react"

export default function ProfilePage() {
  const [user, setUser] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function loadProfile() {
      try {
        const res = await authService.getProfile()
        if (res.success) {
          setUser(res.data)
        }
      } catch (err) {
        console.error("Failed to load profile", err)
      } finally {
        setLoading(false)
      }
    }
    loadProfile()
  }, [])

  if (loading) return <div className="min-h-screen bg-black" />

  return (
    <div className="min-h-screen bg-black pb-20">
      <Navbar />
      <main className="container mx-auto px-6 pt-32">
        <div className="max-w-4xl mx-auto flex flex-col gap-8">
          <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6">
            <div className="flex items-center gap-6">
              <div className="w-24 h-24 md:w-32 md:h-32 rounded-3xl bg-white/5 border border-white/10 flex items-center justify-center">
                <User className="w-12 h-12 md:w-16 md:h-16 text-zinc-500" />
              </div>
              <div>
                <h1 className="text-4xl md:text-5xl font-bold tracking-tighter">{user?.fullName}</h1>
                <div className="flex items-center gap-2 mt-2">
                  <Badge className="bg-white text-black font-bold uppercase tracking-wider px-3">
                    {user?.userRole?.replace("ROLE_", "")}
                  </Badge>
                  {user?.isEmailVerified && (
                    <Badge variant="outline" className="border-emerald-500/50 text-emerald-500 font-bold">
                      VERIFIED
                    </Badge>
                  )}
                </div>
              </div>
            </div>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            <Card className="bg-white/5 border-white/10">
              <CardHeader>
                <CardTitle className="text-xl font-bold">Account Information</CardTitle>
              </CardHeader>
              <CardContent className="flex flex-col gap-4">
                <div className="flex items-center gap-4 text-zinc-400">
                  <Mail className="w-5 h-5" />
                  <div>
                    <p className="text-xs font-bold uppercase tracking-widest text-zinc-600">Email Address</p>
                    <p className="text-white font-medium">{user?.email}</p>
                  </div>
                </div>
                <div className="flex items-center gap-4 text-zinc-400">
                  <Phone className="w-5 h-5" />
                  <div>
                    <p className="text-xs font-bold uppercase tracking-widest text-zinc-600">Phone Number</p>
                    <p className="text-white font-medium">{user?.phoneNumber || "Not provided"}</p>
                  </div>
                </div>
                <div className="flex items-center gap-4 text-zinc-400">
                  <Calendar className="w-5 h-5" />
                  <div>
                    <p className="text-xs font-bold uppercase tracking-widest text-zinc-600">Member Since</p>
                    <p className="text-white font-medium">
                      {new Date(user?.createdAt).toLocaleDateString("en-US", { month: "long", year: "numeric" })}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="bg-white/5 border-white/10">
              <CardHeader>
                <CardTitle className="text-xl font-bold">Reputation & Security</CardTitle>
              </CardHeader>
              <CardContent className="flex flex-col gap-4">
                <div className="flex items-center gap-4 text-zinc-400">
                  <Shield className="w-5 h-5" />
                  <div>
                    <p className="text-xs font-bold uppercase tracking-widest text-zinc-600">Trust Score</p>
                    <p className="text-white font-medium">98/100 (Exceptional)</p>
                  </div>
                </div>
                <div className="p-4 rounded-xl bg-emerald-500/5 border border-emerald-500/10">
                  <p className="text-xs font-bold text-emerald-500 uppercase tracking-widest mb-1">Status</p>
                  <p className="text-sm text-zinc-400">
                    Your account is in good standing. You are eligible to participate in all premium auctions.
                  </p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}
