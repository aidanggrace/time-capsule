"use client"
import { useEffect, useState } from "react"
import { Capsule } from "@/lib/types"
import { supabase } from "@/lib/supabase"

const API_URL = "http://localhost:8080"

export default function Home() {
  const [capsules, setCapsules] = useState<Capsule[]>([])
  const [username, setUsername] = useState("")

  const handleLogout = async () => {
    await supabase.auth.signOut()
    window.location.href = "/login"
  }

  useEffect(() => {
    const getCapsules = async () => {
      const { data: { session } } = await supabase.auth.getSession()
      if (session?.user.email) {
        setUsername(session.user.email.split("@")[0])
      }
      if (!session) {
        window.location.href = "/login"
        return
      }
      const res = await fetch(`${API_URL}/capsules`, {
        headers: { Authorization: `Bearer ${session.access_token}` },
      })

      const data = await res.json()
      setCapsules(data)
    }
    getCapsules()
  }, [])


  return (
    <div>
      <h1>My Capsules</h1>
      <div className="flex items-center justify-between">
        {username && <p>Welcome, {username}</p>}
        <button onClick={handleLogout} className="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600">
          Logout
        </button>
      </div>
      {capsules && capsules.length > 0 ? (
        capsules.map((capsule) => (
          <div key={capsule.id} className="capsule">
            <p>Message: {capsule.message}</p>
            <p>Unlocks: {capsule.unlock_at}</p>
            <p>Sends to: {capsule.recipient_email}</p>
          </div>
        ))
      ) : (
        <p>make a capsule bruh</p>
      )}
    </div>
  )

}

