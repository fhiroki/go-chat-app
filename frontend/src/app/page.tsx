import { FaComments } from "react-icons/fa";

import { cookies } from "next/headers";

export default async function Home() {
  const cookieStore = await cookies();
  const authCookie = cookieStore.get("auth");
  const targetHref = authCookie ? "/chat" : "/login";
  const buttonText = authCookie ? "チャットに入る" : "ログインする";

  return (
    <div className="flex flex-col min-h-screen items-center justify-center bg-gradient-to-br from-blue-500 to-purple-600 p-8">
      <main className="flex flex-col gap-6 items-center">
        <FaComments className="text-white" size={48} />
        <h1 className="text-4xl font-extrabold text-white">チャットアプリへようこそ</h1>
        <p className="text-lg text-center text-white">
          リアルタイムでメッセージを交換し、みんなでコミュニケーションを楽しみましょう。
        </p>
        <div>
          <a
            className="rounded-full bg-white text-blue-600 px-6 py-3 font-semibold shadow hover:bg-gray-100 transition-colors"
            href={targetHref}
          >
            {buttonText}
          </a>
        </div>
      </main>
      <footer className="mt-10">
        <p className="text-sm text-white">© 2025 チャットアプリ</p>
      </footer>
    </div>
  );
}
