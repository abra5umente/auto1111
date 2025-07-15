# Testing the Prompt Generation

## To test the application:

1. **Build the application:**
   ```
   go build -o auto1111.exe .
   ```

2. **Run the application:**
   ```
   ./auto1111.exe
   ```

3. **Test scenarios:**
   - **Empty prompt**: Just press Enter when prompted to test automatic generation
   - **User prompt**: Type something like "a cat" and press Enter to test enhancement

## Expected Debug Output:

When you press Enter (empty prompt), you should see:
```
DEBUG: Starting generate_prompt function
DEBUG: API key found (length: XX)
DEBUG: Creating Gemini client...
DEBUG: Client created successfully
DEBUG: Model configured
DEBUG: Generating new prompt
DEBUG: Calling Gemini API...
DEBUG: API response received, candidates: X
DEBUG: Raw response: "..."
DEBUG: Cleaned result: "..."
Generated prompt: ...
```

## If it hangs:
- The issue is likely in the API call
- Check your internet connection
- Verify your GEMINI_API_KEY is valid
- The API might be rate-limited or down

## If you get an empty prompt:
- The issue is in the response cleaning logic
- Check the "Raw response" debug output to see what Gemini returned