import { render, screen } from '@testing-library/react';
import { cookies } from 'next/headers';
import Home from './page';

// Mock the next/headers module
jest.mock('next/headers', () => ({
  cookies: jest.fn(),
}));

describe('Home component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders login button when not authenticated', async () => {
    // Mock cookies to return no auth cookie
    (cookies as jest.Mock).mockReturnValue({
      get: jest.fn().mockReturnValue(null),
    });

    // @ts-ignore - Render component (ignoring TypeScript error about async component)
    render(await Home());
    
    const heading = screen.getByText('チャットアプリへようこそ');
    expect(heading).toBeInTheDocument();
    
    const loginButton = screen.getByText('ログインする');
    expect(loginButton).toBeInTheDocument();
    expect(loginButton.closest('a')).toHaveAttribute('href', '/login');
  });

  it('renders chat button when authenticated', async () => {
    // Mock cookies to return an auth cookie
    (cookies as jest.Mock).mockReturnValue({
      get: jest.fn().mockReturnValue({ value: 'some-token' }),
    });

    // @ts-ignore - Render component (ignoring TypeScript error about async component)
    render(await Home());
    
    const heading = screen.getByText('チャットアプリへようこそ');
    expect(heading).toBeInTheDocument();
    
    const chatButton = screen.getByText('チャットに入る');
    expect(chatButton).toBeInTheDocument();
    expect(chatButton.closest('a')).toHaveAttribute('href', '/chat');
  });
});