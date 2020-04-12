using System;
using System.Text.Json;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;

namespace netcore
{
    public struct testInput
    {
        public string email { get; set; }
    }

    public struct testOutput
    {
        public int id { get; set; }
        public string name { get; set; }
    }

    public class Startup
    {
        public void Configure(IApplicationBuilder app)
        {
            app.UseMiddleware<JsonMiddleware>();
        }
    }

    public class JsonMiddleware
    {
        private static JsonSerializerOptions JsonSerializerOptions = new JsonSerializerOptions();

        public JsonMiddleware(RequestDelegate next)
        {
        }

        public async Task Invoke(HttpContext context)
        {
            var reqBytes = new byte[(int)context.Request.ContentLength.GetValueOrDefault(32)];

            var readTask = context.Request.Body.ReadAsync(reqBytes, 0, reqBytes.Length);

            var readBytesCount = readTask.IsCompleted ? readTask.Result : await readTask;

            var input = JsonSerializer.Deserialize<testInput>(reqBytes.AsSpan(0, readBytesCount), JsonSerializerOptions);

            var output = new testOutput
            {
                id = int.TryParse(context.Request.Path.Value.AsSpan().Trim('/'), out int t) ? t : 0,
                name = input.email
            };

            context.Response.Headers.Add("Content-Type", "application/json; charset=utf-8");

            var outputBytes = JsonSerializer.SerializeToUtf8Bytes(output, JsonSerializerOptions);

            _ = context.Response.Body.WriteAsync(outputBytes, 0, outputBytes.Length);
        }
    }
}