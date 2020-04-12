using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Hosting;

namespace netcore
{
    public class Program
    {
        public static void Main(string[] args)
        {
            CreateHostBuilder(args).Build().Run();
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>
            Host.CreateDefaultBuilder(args)
            .ConfigureWebHostDefaults(webBuilder =>
            {
                webBuilder.ConfigureLogging(config =>
                {
                    config.ClearProviders();
                });
                webBuilder.ConfigureKestrel(x =>
                {
                    x.AddServerHeader = false;
                    x.AllowSynchronousIO = true;
                });
                webBuilder.UseUrls("http://localhost:5000");
                webBuilder.UseStartup<Startup>();
            });
    }
}