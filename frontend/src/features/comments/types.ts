export interface Comment {
  id: string;
  issueId: string;
  content: string;
  author: string;
  createdAt: string;
  updatedAt?: string;
  mentions?: string[];
}

export interface CommentThreadProps {
  issueId: string;
  className?: string;
}

export interface CommentInputProps {
  onSubmit: (content: string) => void;
  isSubmitting?: boolean;
  placeholder?: string;
}
