package app

import com.typesafe.scalalogging.Logger
import org.joda.time.DateTime
import scopt.OParser

import java.nio.charset.StandardCharsets
import java.nio.file.{Files, Paths}
import scala.io.Source
import scala.util.Using

object Main extends App {
  val logger = Logger("diary-generator")

  val parser = {
    val builder = OParser.builder[Config]
    import builder._
    OParser.sequence(
      opt[String]("base-directory-path").required().action((baseDirectoryPath, config) => config.copy(baseDirectoryPath = baseDirectoryPath)),
      cmd("init")
        .action((_, config) => config.copy(mode = "init"))
        .children(
          opt[String]("template-path").required().action((templatePath, config) => config.copy(templatePathOpt = Some(templatePath)))
        ),
      cmd("archive")
        .action((_, config) => config.copy(mode = "archive"))
        .children(
          opt[String]("starts-with").required().action((startsWith, config) => config.copy(startsWithOpt = Some(startsWith)))
        )
    )
  }

  def templateFormatter(template: String): String = {
    val now = DateTime.now()
    template.replaceAll("%TODAY%", Format.ymd(now))
  }

  def init(baseDirectoryPath: String, templateFilePath: String, dt: DateTime): Unit = {
    val baseDir = Paths.get(baseDirectoryPath)

    val templateFile = Paths.get(templateFilePath).toFile

    val targetDirPath = baseDir.resolve(Format.ymd(dt)).toAbsolutePath
    val targetDir = targetDirPath.toFile
    val targetFilePath = targetDirPath.resolve(templateFile.getName)
    val targetFile = targetFilePath.toFile

    if (templateFile.exists()) {
      if (targetDir.exists()) {
        logger.info(s"already exist : ${targetDirPath.toString}")
      } else {
        logger.info(s"make directory : ${targetDirPath.toString}")
        targetDir.mkdirs()
      }

      if (targetFile.exists()) {
        logger.info(s"already exist : ${targetFilePath.toString}")
      } else {
        logger.info(s"create file : ${targetFilePath.toString}")

        Using(Source.fromFile(templateFile, "UTF-8", 4096)) { source => source.mkString }.fold(
          err => {
            logger.error("can't open template file", err)
          },
          template => {
            val formattedTemplate = templateFormatter(template)
            targetFile.createNewFile()
            Files.write(targetFilePath, formattedTemplate.getBytes(StandardCharsets.UTF_8))
          }
        )
      }
    } else {
      logger.error("template file not found")
    }
  }

  OParser
    .parse(parser, args, Config())
    .fold {
      // do nothing
    } { config =>
      (config.mode, config.templatePathOpt, config.startsWithOpt) match {
        case ("init", Some(templateFilePath), _) =>
          init(config.baseDirectoryPath, templateFilePath, DateTime.now())
        case ("archive", _, Some(_)) => ??? // todo impl
        case _ => ??? // todo impl
      }
    }
}
