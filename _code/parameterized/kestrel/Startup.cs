using System;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;

namespace netcore
{
    public class Startup
    {
        public void Configure(IApplicationBuilder app)
        {
            app.UseMiddleware<ParamMiddleware>();
        }
    }

    public class ParamMiddleware
    {
        private static byte[] BodyBytes = System.Text.Encoding.UTF8.GetBytes("Hello ");

        public ParamMiddleware(RequestDelegate next)
        {
        }

        public Task Invoke(HttpContext context)
        {
            var nameSpan = context.Request.Path.Value.AsSpan().TrimStart("/hello/");
            var outputBytes = new byte[BodyBytes.Length + nameSpan.Length];
            BodyBytes.CopyTo(outputBytes, 0);
            var outputBytesCount = System.Text.Encoding.UTF8.GetBytes(nameSpan, outputBytes.AsSpan(BodyBytes.Length));

            return context.Response.Body.WriteAsync(outputBytes, 0, BodyBytes.Length + outputBytesCount);
        }
    }
}