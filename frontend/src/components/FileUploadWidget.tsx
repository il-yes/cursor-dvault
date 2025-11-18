import { useState, useRef } from "react";
import { Upload, X, FileText, File } from "lucide-react";
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
    <div className="space-y-3">
      {/* Upload Area */}
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => fileInputRef.current?.click()}
        className={cn(
          "relative border-2 border-dashed rounded-lg p-6 transition-all cursor-pointer",
          "hover:border-primary/50 hover:bg-accent/5",
          isDragging
            ? "border-primary bg-accent/10"
            : "border-border bg-background"
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

        <div className="flex flex-col items-center justify-center space-y-2 text-center">
          <div className="p-3 rounded-full bg-accent/10">
            <Upload className="w-6 h-6 text-muted-foreground" />
          </div>
          <div className="space-y-1">
            <p className="text-sm font-medium text-foreground">
              Drop files here or click to browse
            </p>
            <p className="text-xs text-muted-foreground">
              Supports: {acceptedTypes.join(", ")} (max {maxFiles} files)
            </p>
          </div>
        </div>
      </div>

      {/* Selected Files List */}
      {selectedFiles.length > 0 && (
        <div className="space-y-2">
          <p className="text-xs font-medium text-muted-foreground">
            Selected Files ({selectedFiles.length}/{maxFiles})
          </p>
          <div className="space-y-2">
            {selectedFiles.map((file, index) => (
              <div
                key={index}
                className="flex items-center gap-3 p-3 rounded-lg bg-accent/5 border border-border group hover:border-primary/30 transition-colors"
              >
                <div className="flex-shrink-0 p-2 rounded bg-primary/10">
                  {file.type.startsWith("image/") ? (
                    <FileText className="w-4 h-4 text-primary" />
                  ) : (
                    <File className="w-4 h-4 text-primary" />
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-foreground truncate">
                    {file.name}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {formatFileSize(file.size)}
                  </p>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={(e) => {
                    e.stopPropagation();
                    removeFile(index);
                  }}
                  className="opacity-0 group-hover:opacity-100 transition-opacity"
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