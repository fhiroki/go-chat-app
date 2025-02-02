"use client";
import Head from "next/head";
import Link from "next/link";
import { FaGoogle, FaGithub } from "react-icons/fa";

const Login = () => {
  return (
    <>
      <Head>
        <title>Login - My App</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </Head>
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-md">
          <h1 className="text-3xl font-bold text-center">Please login</h1>
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            <h3 className="text-xl font-semibold">Sign in required for chat</h3>
            <p className="mt-2">Please select a provider to sign in:</p>
          </div>
          <div className="space-y-4">
            <Link href="http://localhost:8080/auth/login/google" className="flex items-center justify-center gap-2 px-4 py-2 border border-transparent rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 transition">
              <FaGoogle size={20} />
              Google
            </Link>
            <Link href="http://localhost:8080/auth/login/github" className="flex items-center justify-center gap-2 px-4 py-2 border border-transparent rounded-md shadow-sm text-white bg-gray-800 hover:bg-gray-900 transition">
              <FaGithub size={20} />
              GitHub
            </Link>
          </div>
        </div>
      </div>
    </>
  );
};

export default Login;
