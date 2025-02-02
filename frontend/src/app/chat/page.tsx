"use client";
import { useEffect, useState, useRef } from "react";
import axios from "axios";
import Head from "next/head";
import Image from "next/image";

interface MessageType {
  id: number;
  email: string;
  avatarURL: string;
  message: string;
  createdAt: string;
}

export default function ChatPage() {
  const [messages, setMessages] = useState<MessageType[]>([]);
  const [input, setInput] = useState("");
  const socketRef = useRef<WebSocket | null>(null);

  // ユーザのemailをクッキー（auth）から取得する関数
  const getUserEmail = () => {
    const cookies = document.cookie.split("; ");
    for (const cookie of cookies) {
      if (cookie.startsWith("auth=")) {
        const encodedValue = cookie.substring("auth=".length);
        try {
          const decoded = JSON.parse(atob(encodedValue));
          return decoded.email;
        } catch (error) {
          console.error("Failed to get user email:", error);
        }
      }
    }
    return "unknown@example.com";
  };

  const userEmail = getUserEmail();

  useEffect(() => {
    if (!window.WebSocket) {
      alert("error: WebSocket not supported");
      return;
    }
    socketRef.current = new WebSocket("ws://localhost:8080/room");

    socketRef.current.onmessage = (e) => {
      try {
        const msgData = JSON.parse(e.data);
        let createdAt = "";
        if (msgData.CreatedAt) {
          const t = new Date(msgData.CreatedAt);
          createdAt = t.toLocaleString();
        }
        const newMessage: MessageType = {
          id: msgData.ID,
          email: msgData.Email,
          avatarURL: msgData.AvatarURL,
          message: msgData.Message,
          createdAt: createdAt,
        };
        setMessages((prev) => [...prev, newMessage]);
      } catch (error) {
        console.error("Error parsing message:", error);
      }
    };

    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, []);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!input.trim()) return;
    if (!socketRef.current || socketRef.current.readyState !== WebSocket.OPEN) {
      alert("error: socket not connected or still connecting");
      return;
    }
    socketRef.current.send(JSON.stringify({ Message: input }));
    setInput("");
  };

  // メッセージを削除するハンドラー
  const handleDelete = (id: number) => {
    axios
      .delete(`http://localhost:8080/messages/delete?id=${id}`)
      .then(() => {
        setMessages((prev) => prev.filter((msg) => msg.id !== id));
      })
      .catch((error) => {
        console.error("Error deleting message:", error);
      });
  };

  return (
    <>
      <Head>
        <title>Chat</title>
      </Head>
      <div className="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center p-6">
        <div className="bg-white rounded-2xl shadow-2xl flex flex-col w-full max-w-2xl h-full max-h-[90vh] overflow-hidden">
          {/* ヘッダー */}
          <div className="bg-gray-50 border-b border-gray-100 py-4 px-6 flex justify-between items-center shadow-sm">
            <div className="flex items-center space-x-4">
              <span className="font-medium text-gray-800">{userEmail}</span>
              <a
                href="http://localhost:8080/logout"
                className="text-sm text-blue-600 hover:text-blue-800 transition-colors"
              >
                Logout
              </a>
            </div>
            <h3 className="text-2xl font-semibold text-gray-800">Chat Room</h3>
          </div>
          {/* チャットエリア */}
          <div className="flex-1 overflow-y-auto p-6">
            <ul id="messages" className="space-y-4">
              {messages.map((msg, index) => (
                <li key={index} className="flex items-start space-x-4">
                  <Image
                    title={msg.email}
                    src={msg.avatarURL}
                    alt="avatar"
                    width={40}
                    height={40}
                    className="rounded-full object-cover"
                  />
                  <div className="flex flex-col w-full">
                    <div className="flex items-center">
                      <span className="font-semibold text-gray-800">
                        {msg.email}
                      </span>
                      <span className="text-xs text-gray-500 ml-2">
                        [{msg.createdAt}]
                      </span>
                      {userEmail === msg.email && (
                        <button
                          onClick={() => handleDelete(msg.id)}
                          className="ml-auto text-xs text-red-500 hover:text-red-600 transition-colors"
                        >
                          削除
                        </button>
                      )}
                    </div>
                    <div className="mt-2 bg-gray-100 p-4 rounded-xl shadow transition duration-300 ease-in-out hover:shadow-md">
                      <span className="text-gray-700">{msg.message}</span>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
          {/* チャットフォーム */}
          <form
            id="chatbox"
            role="form"
            onSubmit={handleSubmit}
            className="border-t border-gray-200 bg-gray-50 p-4 flex items-center"
          >
            <textarea
              id="message"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Type your message..."
              rows={2}
              className="text-gray-800 flex-1 border border-gray-300 rounded-md p-2 mr-2 focus:outline-none focus:ring-2 focus:ring-blue-400 transition-colors"
            ></textarea>
            <button
              type="submit"
              className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded"
            >
              Send
            </button>
          </form>
        </div>
      </div>
    </>
  );
}
