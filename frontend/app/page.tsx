"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";

import { getCookie } from "@/utils/cookies";
import { backendBase } from "@/utils/util";

export default function Home() {
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    useEffect(() => {
        const token = getCookie("ppet_token");
        if (token) {
            setIsLoggedIn(true);
        }
    }, []);

    const handleLogin = () => {
        window.location.href = `${backendBase}/auth/google/login`;
    };

    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-24">
            {isLoggedIn ? (
                <h1 className="mb-4 text-2xl">
                    Welcome back! You are logged in.
                </h1>
            ) : (
                <>
                    <h1 className="mb-4 text-2xl">
                        Gochi. Sign in with Google to continue.
                    </h1>
                    <Button onClick={handleLogin}>Sign in with Google</Button>
                </>
            )}
        </main>
    );
}
