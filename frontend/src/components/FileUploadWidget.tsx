import { useState, useRef } from "react";
import { Upload, X, FileText, File, Badge, Image } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface FileUploadWidgetProps {
  onFileSelect: (files: File[]) => void;
  maxFiles?: number;
  acceptedTypes?: string[];
  value?: File[];
}

const ACCEPTED_TYPES = [
  ".jpeg",
  ".jpg",
  ".svg",
  ".png",
  ".pdf",
  ".doc",
  ".docx",
];

export function FileUploadWidget({
  onFileSelect,
  maxFiles = 5,
  acceptedTypes = ACCEPTED_TYPES,
  value = [],
}: FileUploadWidgetProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [selectedFiles, setSelectedFiles] = useState<File[]>(value);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    handleFiles(files);
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files);
      handleFiles(files);
    }
  };

  const handleFiles = (files: File[]) => {
    // Filter by accepted types
    const validFiles = files.filter((file) => {
      const ext = `.${file.name.split(".").pop()?.toLowerCase()}`;
      return acceptedTypes.includes(ext);
    });

    // Respect max files limit
    const newFiles = [...selectedFiles, ...validFiles].slice(0, maxFiles);
    setSelectedFiles(newFiles);
    onFileSelect(newFiles);
  };

  const removeFile = (index: number) => {
    const newFiles = selectedFiles.filter((_, i) => i !== index);
    setSelectedFiles(newFiles);
    onFileSelect(newFiles);
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + " " + sizes[i];
  };

  return (
  <div className="space-y-4">
    {/* Glass Upload Area */}
    <div
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      onClick={() => fileInputRef.current?.click()}
      className={cn(
        "group relative rounded-3xl p-10 transition-all cursor-pointer backdrop-blur-xl",
        "border-2 border-white/30 bg-white/50 dark:bg-zinc-900/50 shadow-2xl hover:shadow-primary/20 hover:shadow-3xl",
        "hover:border-primary/50 hover:bg-white/70 dark:hover:bg-zinc-900/70",
        isDragging 
          ? "border-primary/70 bg-gradient-to-br from-primary/20 to-primary/10 shadow-primary/30 scale-[1.02]" 
          : "border-white/40 bg-gradient-to-br from-white/60 via-white/40 to-zinc-100/30 dark:from-zinc-900/60 dark:via-zinc-900/40 dark:to-black/30"
      )}
    >
      <input
        ref={fileInputRef}
        type="file"
        multiple
        accept={acceptedTypes.join(",")}
        onChange={handleFileInput}
        className="hidden"
      />

      {/* Animated glow overlay */}
      <div className="absolute inset-0 rounded-3xl bg-gradient-to-r from-primary via-primary/30 to-primary opacity-0 group-hover:opacity-20 transition-opacity blur-xl" />
      
      <div className="relative flex flex-col items-center justify-center space-y-4 text-center">
        <div className="relative group/icon">
          <div className="absolute inset-0 bg-gradient-to-br from-primary/30 to-amber-500/30 rounded-3xl blur-xl opacity-0 group-hover/icon:opacity-60 transition-all" />
          <div className="p-4 rounded-3xl bg-white/60 dark:bg-zinc-800/60 backdrop-blur-xl border border-white/50 shadow-xl group-hover/icon:shadow-2xl transition-all">
            <Upload className="w-8 h-8 text-primary drop-shadow-lg" />
          </div>
        </div>
        <div className="space-y-1">
          <p className="text-lg font-bold bg-gradient-to-r from-foreground to-primary/80 bg-clip-text text-transparent">
            Drop files here or click to browse
          </p>
          <p className="text-sm text-muted-foreground/90 font-medium">
            Supports: {acceptedTypes.slice(0, 4).join(", ")}{acceptedTypes.length > 4 && "..."} (max {maxFiles} files)
          </p>
        </div>
      </div>
    </div>

    {/* Selected Files Glass List */}
    {selectedFiles.length > 0 && (
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <p className="text-sm font-semibold bg-gradient-to-r from-foreground to-muted-foreground bg-clip-text text-transparent">
            Selected Files ({selectedFiles.length}/{maxFiles})
          </p>
          {selectedFiles.length > 0 && (
            <Badge className="bg-gradient-to-r from-primary/20 to-amber-500/20 backdrop-blur-sm border-primary/30 text-primary font-semibold px-3 py-1 shadow-md">
              {selectedFiles.length} file{selectedFiles.length !== 1 ? 's' : ''}
            </Badge>
          )}
        </div>
        
        <div className="space-y-3 max-h-48 overflow-y-auto glass-scrollbar">
          {selectedFiles.map((file, index) => (
            <div
              key={index}
              className="group flex items-center gap-4 p-4 rounded-2xl backdrop-blur-xl bg-white/70 dark:bg-zinc-900/70 border border-white/40 shadow-xl hover:shadow-2xl hover:shadow-primary/20 hover:scale-[1.01] transition-all duration-300 hover:border-primary/50"
            >
              {/* File icon */}
              <div className="flex-shrink-0 p-3 rounded-2xl bg-gradient-to-br from-primary/20 to-amber-500/20 backdrop-blur-sm border border-primary/30 shadow-md">
                {file.type.startsWith("image/") ? (
                  <Image className="w-5 h-5 text-primary" />
                ) : (
                  <FileText className="w-5 h-5 text-primary" />
                )}
              </div>
              
              {/* File info */}
              <div className="flex-1 min-w-0">
                <p className="text-sm font-semibold text-foreground truncate group-hover:text-primary transition-colors">
                  {file.name}
                </p>
                <p className="text-xs text-muted-foreground/80 font-mono">
                  {formatFileSize(file.size)}
                </p>
              </div>
              
              {/* Remove button */}
              <Button
                type="button"
                variant="ghost"
                size="icon"
                onClick={(e) => {
                  e.stopPropagation();
                  removeFile(index);
                }}
                className="h-10 w-10 rounded-2xl bg-white/80 dark:bg-zinc-800/80 border border-destructive/30 opacity-70 hover:opacity-100 hover:bg-destructive/10 hover:border-destructive/50 hover:text-destructive shadow-sm transition-all group-hover:opacity-100"
              >
                <X className="w-4 h-4" />
              </Button>
            </div>
          ))}
        </div>
      </div>
    )}
  </div>
);

}