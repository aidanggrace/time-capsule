"use client"
import { useEffect, useState } from "react"
import { Capsule } from "@/lib/types"

const API_URL = "http://localhost:8080"

export default function Home() {
  const [capsules, setCapsules] = useState<Capsule[]>([])

  useEffect(() => {
    const getCapsules = async () => {
      const res = await fetch(`${API_URL}/capsules`, { credentials: "include" })
      if (res.status === 401) {
        window.location.href = "/login"
        return
      }
      const data = await res.json()
      setCapsules(data)
    }
    getCapsules()
  }, [])


  return (
    <div>
      <h1>My Capsules</h1>
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

