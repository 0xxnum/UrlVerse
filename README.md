# UrlVerse

**UrlVerse** helps you clean up your list of URLs, making it useful for pentesting and bug bounty.

## What It Does

This tool filters out less interesting URLs from your list, helping you focus on what really matters. For instance, if you’re digging through results from tools like Waymore, UrlVerse can give you a fresh start by filtering out the noise.

## Equivalences

The **equivalences** section allows you to define specific words or terms that should be treated as equivalent during processing. This is useful when you want to group different variations of a term together to ensure they are recognized as the same.

### Example

If you're interested in processing URLs related to Tesla vehicles, you might define equivalences like this:

    equivalences := map[string][]string{
    "TESLA": {"model-3", "model-y", "model-s", "model-x"},
        }
## Key Features

- **Excludes boring URLs**: UrlVerse prints one URL per feature of a website while blocking known dull URLs.
- **Customizable filters**: Easily add your own paths and extensions to the exclusion list.
- **Smart detection**: It can  recognize patterns in URLs, ignoring things like:
  - Directories such as `/docs`
  - Files with extensions like `.png`
  - Profile pages
  - Certain query parameters

## What It Ignores

UrlVerse is designed to filter out:
- **Boring directories**: e.g., `/docs`, `/support`
- **Uninteresting file types**: e.g., `.css`, `.jpg`
- **Profile and blog pages**: e.g., `/user/FooBar`
- **Certain URL parameters**: e.g., UTM parameters

## Usage

To get started with UrlVerse, follow these steps:

### Installation

## Downloading and Running UrlVerse

To get started with UrlVerse, follow these steps:

### Clone the Repository
Make sure you have Git installed, then run the following command in your terminal:
 
    git clone https://github.com/0xxnum/UrlVerse.git
   
    cd 0xxnum-UrlVerse

 Run UrlVerse
    You can run UrlVerse by piping your URLs into it:
    
    cat many_urls.txt | go run urlverse.py | tee less_urls.txtx
# Practical Example
For a practical example, you can filter URLs from Waymore like this:

    waymore example.org | tee all_urls.txt | python urlverse.py > filtered_urls.txt


# Contributing
If you have ideas for additional filters or find a bug, feel free to reach out! Your contributions are always welcome.
